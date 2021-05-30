package facebook

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tusa-plus/core/utils"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Facebook struct {
	logger         *zap.Logger
	httpClientPool *utils.HttpClientPool
}

const (
	fbURL = "https://graph.facebook.com/me?"
)

var ErrDoRequest = fmt.Errorf("failed to request")
var ErrValidate = fmt.Errorf("failed to validate result")

func (fb *Facebook) GetEmail(ctx context.Context, fbToken string) (string, error) {
	params := url.Values{}
	params.Add("fields", "email")
	params.Add("access_token", fbToken)
	request, err := http.NewRequest("GET", fbURL+params.Encode(), nil)
	if err != nil {
		fb.logger.Error("unexpected error during creating request",
			zap.Error(err),
			zap.String("access_token", fbToken),
		)
		return "", ErrDoRequest
	}
	client := fb.httpClientPool.Get()
	defer fb.httpClientPool.Put(client)
	response, err := client.Do(request.WithContext(ctx))
	if err != nil {
		fb.logger.Error("unexpected error during request",
			zap.Error(err),
			zap.String("access_token", fbToken),
		)
		return "", ErrDoRequest
	}
	defer func() {
		if err := response.Body.Close(); err != nil { //nolint:govet
			fb.logger.Error("unexpected error during body close",
				zap.Error(err),
			)
		}
	}()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fb.logger.Error("unexpected error during read body",
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
