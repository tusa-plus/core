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

const TokenTypeAccess = "access"
const TokenTypeRefresh = "refresh"

const TokenTypeProperty = "token_type"
const TokenIdProperty = "token_id"
const TokenExpProperty = "exp"

var ErrTokenExpired = fmt.Errorf("token expired")
var ErrInvalidSignature = fmt.Errorf("token has invalid signature")
var ErrInvalidToken = fmt.Errorf("token is invalid")
var ErrExpireNonRefresh = fmt.Errorf("use refresh token to expire pair")

func (ts *TokenStorage) NewTokenPair(data map[string]interface{}) (string, string, error) {
	tokenId := uuid.NewV4().String()
	claims := jwt.MapClaims{
		TokenIdProperty: tokenId,
	}
	for key, value := range data {
		claims[key] = value
	}
	// create access token
	claims[TokenExpProperty] = time.Now().Add(ts.accessExpiration).Unix()
	claims[TokenTypeProperty] = TokenTypeAccess
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := unsignedToken.SignedString(ts.secret)
	if err != nil {
		return "", "", nil
	}
	// create refresh token
	claims[TokenExpProperty] = time.Now().Add(ts.refreshExpiration).Unix()
	claims[TokenTypeProperty] = TokenTypeRefresh
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
			return nil, ErrInvalidToken
		}
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid || claims.Valid() != nil {
		return nil, ErrInvalidToken
	}
	tokenIdRaw, ok := claims[TokenIdProperty]
	if !ok {
		return nil, ErrInvalidToken
	}
	tokenId, ok := tokenIdRaw.(string)
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
	tokenTypeRaw, ok := token[TokenTypeProperty]
	if !ok {
		return ErrInvalidToken
	}
	tokenType, ok := tokenTypeRaw.(string)
	if !ok || tokenType != TokenTypeRefresh {
		return ErrExpireNonRefresh
	}
	expAt, err := token[TokenExpProperty].(json.Number).Int64()
	if err != nil {
		return ErrInvalidToken
	}
	tokenId := token[TokenIdProperty].(string)
	return ts.storage.Set(tokenId, []byte{0}, time.Unix(expAt, 0).Add(time.Second).Sub(time.Now()))
}
