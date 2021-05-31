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

type Google struct {
	logger         *zap.Logger
	httpClientPool *utils.HTTPClientPool
}

const (
	googleURL = "https://www.googleapis.com/userinfo/v2/me"
)

var ErrDoRequest = fmt.Errorf("failed to request")
var ErrValidate = fmt.Errorf("failed to validate result")

func (google *Google) GetEmail(ctx context.Context, googleToken string) (string, error) {
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
