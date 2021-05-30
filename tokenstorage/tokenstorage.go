package tokenstorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/twinj/uuid"
	"go.uber.org/zap"
	"time"
)

type TokenStorage struct {
	logger            *zap.Logger
	secret            []byte
	storage           fiber.Storage
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

const TokenTypeAccess = "access"
const TokenTypeRefresh = "refresh"

const TokenTypeProperty = "token_type"
const tokenIDProperty = "token_id"
const TokenExpProperty = "exp"

var ErrTokenExpired = fmt.Errorf("token expired")
var ErrInvalidSignature = fmt.Errorf("token has invalid signature")
var ErrInvalidToken = fmt.Errorf("token is invalid")
var ErrExpireNonRefresh = fmt.Errorf("use refresh token to expire pair")
var ErrMissingStorage = fmt.Errorf("missing storage in create")
var ErrUnknown = fmt.Errorf("unknow error")

func NewTokenStorage(
	secret []byte,
	logger *zap.Logger,
	storage fiber.Storage,
	accessExpiration time.Duration,
	refreshExpiration time.Duration,
) (*TokenStorage, error) {
	if logger == nil {
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("failed to create zap logger: %w", err)
		}
	}
	if storage == nil {
		return nil, ErrMissingStorage
	}
	return &TokenStorage{
		logger:            logger,
		secret:            secret,
		storage:           storage,
		accessExpiration:  accessExpiration,
		refreshExpiration: refreshExpiration,
	}, nil
}

func (ts *TokenStorage) NewTokenPair(data map[string]interface{}) (string, string, error) {
	tokenID := uuid.NewV4().String()
	claims := jwt.MapClaims{
		tokenIDProperty: tokenID,
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
		ts.logger.Error("unexpected error creating token",
			zap.Error(err),
		)
		return "", "", ErrUnknown
	}
	// create refresh token
	claims[TokenExpProperty] = time.Now().Add(ts.refreshExpiration).Unix()
	claims[TokenTypeProperty] = TokenTypeRefresh
	unsignedToken = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := unsignedToken.SignedString(ts.secret)
	if err != nil {
		ts.logger.Error("unexpected error creating token",
			zap.Error(err),
		)
		return "", "", ErrUnknown
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
		var validationError *jwt.ValidationError
		if errors.As(err, &validationError) {
			if (validationError.Errors & jwt.ValidationErrorExpired) > 0 {
				return nil, ErrTokenExpired
			}
			if (validationError.Errors & (jwt.ValidationErrorSignatureInvalid)) > 0 {
				return nil, ErrInvalidSignature
			}
			return nil, ErrInvalidToken
		}
		ts.logger.Error("unexpected error during parse token",
			zap.Error(err),
			zap.String("token_string", tokenString),
		)
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid || claims.Valid() != nil {
		return nil, ErrInvalidToken
	}
	tokenIDRaw, ok := claims[tokenIDProperty]
	if !ok {
		ts.logger.Error("token doesn't contain token_id",
			zap.String("token_string", tokenString),
		)
		return nil, ErrUnknown
	}
	tokenID, ok := tokenIDRaw.(string)
	if !ok {
		ts.logger.Error("token_id is not string",
			zap.String("token_string", tokenString),
		)
		return nil, ErrUnknown
	}
	data, err := ts.storage.Get(tokenID)
	if data == nil && err == nil {
		return claims, nil
	}
	if err != nil {
		ts.logger.Error("unexpected error during storage get",
			zap.Error(err),
			zap.String("token_string", tokenString),
			zap.String("token_id", tokenID),
		)
		return nil, ErrUnknown
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
		ts.logger.Error("token doesn't contain token_type",
			zap.Error(err),
			zap.String("token_string", tokenString),
		)
		return ErrUnknown
	}
	tokenType, ok := tokenTypeRaw.(string)
	if !ok {
		ts.logger.Error("token_type is not string",
			zap.String("token_string", tokenString),
		)
		return ErrUnknown
	}
	if tokenType != TokenTypeRefresh {
		return ErrExpireNonRefresh
	}
	expAtRaw, ok := token[TokenExpProperty]
	if !ok {
		ts.logger.Error("token doesn't contain exp",
			zap.String("token_string", tokenString),
		)
		return ErrUnknown
	}
	expAt, err := expAtRaw.(json.Number).Int64()
	if err != nil {
		return ErrUnknown
	}
	tokenIDRaw, ok := token[tokenIDProperty]
	if !ok {
		ts.logger.Error("token doesn't contain token_id",
			zap.String("token_string", tokenString),
		)
		return ErrUnknown
	}
	tokenID, ok := tokenIDRaw.(string)
	if !ok {
		ts.logger.Error("token_id is not string",
			zap.String("token_string", tokenString),
		)
		return ErrUnknown
	}
	err = ts.storage.Set(tokenID, []byte{0}, time.Until(time.Unix(expAt, 0).Add(time.Second)))
	if err != nil {
		ts.logger.Error("unexpected error set token_id",
			zap.Error(err),
			zap.String("token_string", tokenString),
		)
		return ErrUnknown
	}
	return nil
}
