package yandexgames

import (
	"fmt"

	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

	"github.com/tusa-plus/core/utils"
	"go.uber.org/zap"
)

var ErrValidate = fmt.Errorf("failed to validate result")
var ErrInvalidData = fmt.Errorf("invalid data")

type YandexGames struct {
	logger         *zap.Logger
	httpClientPool *utils.HTTPClientPool
	secret         []byte
}

func (yg *YandexGames) ValidateSignature(sign string, data string) error {
	message, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		yg.logger.Error("unexpected error during decoding data",
			zap.Error(err),
			zap.String("data", data),
		)
		return ErrInvalidData
	}
	h := hmac.New(sha256.New, yg.secret)
	h.Write(message)
	result := base64.StdEncoding.EncodeToString(h.Sum(nil))
	if result != sign {
		return ErrValidate
	}
	return nil
}
