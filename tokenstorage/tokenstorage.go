package tokenstorage

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/twinj/uuid"
	"time"
)

type TokenStorage struct {
	secret            []byte
	storage           fiber.Storage
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

var ErrTokenExpired = fmt.Errorf("token expired")
var ErrInvalidSignature = fmt.Errorf("token has invalid signature")
var ErrInvalidToken = fmt.Errorf("token is invalid")

func (ts *TokenStorage) NewTokenPair(data map[string]interface{}) (string, string, error) {
	tokenId := uuid.NewV4().String()
	claims := jwt.MapClaims{
		"token_id": tokenId,
	}
	for key, value := range data {
		claims[key] = value
	}
	// create access token
	claims["exp"] = time.Now().Add(ts.accessExpiration).Unix()
	claims["token_type"] = "access"
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := unsignedToken.SignedString(ts.secret)
	if err != nil {
		return "", "", nil
	}
	// create refresh token
	claims["exp"] = time.Now().Add(ts.refreshExpiration).Unix()
	claims["token_type"] = "refresh"
	unsignedToken = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := unsignedToken.SignedString(ts.secret)
	if err != nil {
		return "", "", nil
	}
	return accessToken, refreshToken, nil
}

func (ts *TokenStorage) ParseToken(tokenString string) (map[string]interface{}, error) {
	parser := jwt.Parser{
		UseJSONNumber:        true,
		SkipClaimsValidation: false,
	}
	token, err := parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSignature
		}
		return ts.secret, nil
	})
	if err != nil {
		if validationError, ok := err.(*jwt.ValidationError); ok {
			if (validationError.Errors & jwt.ValidationErrorExpired) > 0 {
				return nil, ErrTokenExpired
			}
			if (validationError.Errors & (jwt.ValidationErrorSignatureInvalid)) > 0 {
				return nil, ErrInvalidSignature
			}
		}
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid || claims.Valid() != nil {
		return nil, ErrInvalidToken
	}
	tokenId, ok := claims["token_id"].(string)
	if !ok {
		return nil, ErrInvalidToken
	}
	data, err := ts.storage.Get(tokenId)
	if data == nil && err == nil {
		return claims, nil
	}
	if err != nil {
		return nil, err
	}
	return nil, ErrTokenExpired
}

func (ts *TokenStorage) ExpireToken(tokenString string) error {
	token, err := ts.ParseToken(tokenString)
	if err != nil {
		return err
	}
	expAt, err := token["exp"].(json.Number).Int64()
	if err != nil {
		return ErrInvalidToken
	}
	tokenId := token["token_id"].(string)
	return ts.storage.Set(tokenId, []byte{0}, time.Unix(expAt, 0).Add(time.Second).Sub(time.Now()))
}
