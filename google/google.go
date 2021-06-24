package google

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tusa-plus/core/utils"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type Google interface {
	GetEmail(ctx context.Context, googleToken string) (string, error)
}

type googleDefaultImpl struct {
	logger         *zap.Logger
	httpClientPool *utils.HTTPClientPool
}

type googleMock struct{}

func NewMockGoogle() Google {
	return &googleMock{}
}

func NewGoogle(logger *zap.Logger, pool *utils.HTTPClientPool) (Google, error) {
	if logger == nil {
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("failed to create zap logger: %w", err)
		}
	}
	if pool == nil {
		newPool := utils.NewHTTPClientPool()
		pool = &newPool
	}
	google := &googleDefaultImpl{
		logger:         logger,
		httpClientPool: pool,
	}
	return google, nil
}

const (
	googleURL = "https://www.googleapis.com/userinfo/v2/me"
)

var ErrDoRequest = fmt.Errorf("failed to request")
var ErrValidate = fmt.Errorf("failed to validate result")

func (google *googleDefaultImpl) GetEmail(ctx context.Context, googleToken string) (string, error) {
	request, err := http.NewRequest("GET", googleURL, nil)
	if err != nil {
		google.logger.Error("unexpected error during creating request",
			zap.Error(err),
			zap.String("access_token", googleToken),
		)
		return "", ErrDoRequest
	}
	request.Header.Add("Authorization", "Bearer "+googleToken)
	client := google.httpClientPool.Get()
	defer google.httpClientPool.Put(client)
	response, err := client.Do(request.WithContext(ctx))
	if err != nil {
		google.logger.Error("unexpected error during request",
			zap.Error(err),
			zap.String("access_token", googleToken),
		)
		return "", ErrDoRequest
	}
	defer func() {
		if err = response.Body.Close(); err != nil {
			google.logger.Error("unexpected error during body close",
				zap.Error(err),
			)
		}
	}()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		google.logger.Error("unexpected error during read body",
			zap.Error(err),
		)
		return "", ErrDoRequest
	}
	var responseJSON map[string]json.RawMessage
	if err := json.Unmarshal(body, &responseJSON); err != nil {
		return "", ErrValidate
	}
	var email string
	if err := json.Unmarshal(responseJSON["email"], &email); err != nil {
		return "", ErrValidate
	}
	return email, nil
}

func (fb *googleMock) GetEmail(_ context.Context, fbToken string) (string, error) {
	return fbToken, nil
}
