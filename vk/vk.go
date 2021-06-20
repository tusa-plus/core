package vk

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

type Vk struct {
	logger         *zap.Logger
	httpClientPool *utils.HTTPClientPool
}

const (
	vkURL = "https://api.vk.com/method/users.get?"
)

var ErrDoRequest = fmt.Errorf("failed to request")
var ErrValidate = fmt.Errorf("failed to validate result")

///todo ничего пока не работает
func (vk *Vk) GetID(ctx context.Context, vkToken string) (uint64, error) {
	params := url.Values{}
	params.Add("access_token", vkToken)
	params.Add("v", "5.131")
	request, err := http.NewRequest("GET", vkURL+params.Encode(), nil)
	if err != nil {
		vk.logger.Error("unexpected error during creating request",
			zap.Error(err),
			zap.String("access_token", vkToken),
		)
		return 0, ErrDoRequest
	}
	client := vk.httpClientPool.Get()
	defer vk.httpClientPool.Put(client)
	response, err := client.Do(request.WithContext(ctx))
	if err != nil {
		vk.logger.Error("unexpected error during request",
			zap.Error(err),
			zap.String("access_token", vkToken),
		)
		return 0, ErrDoRequest
	}
	defer func() {
		if err = response.Body.Close(); err != nil {
			vk.logger.Error("unexpected error during body close",
				zap.Error(err),
			)
		}
	}()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		vk.logger.Error("unexpected error during read body",
			zap.Error(err),
		)
		return 0, ErrDoRequest
	}
	var responseJSON map[string]json.RawMessage
	if err := json.Unmarshal(body, &responseJSON); err != nil {
		return 0, ErrValidate
	}
	var answer []json.RawMessage
	if err := json.Unmarshal(responseJSON["response"], &answer); err != nil {
		return 0, ErrValidate
	}
	var profileInfo map[string]json.RawMessage
	if err := json.Unmarshal(answer[0], &profileInfo); err != nil {
		return 0, ErrValidate
	}
	var id uint64
	if err := json.Unmarshal(profileInfo["id"], &id); err != nil {
		return 0, ErrValidate
	}
	return id, nil
}
