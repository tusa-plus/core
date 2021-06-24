package vk

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"time"
)

type Vk interface {
	GetAccount(vkToken string) (*Account, error)
}

type Account struct {
	Id      uint64 `json:"id"`
	Name    string `json:"first_name"`
	Surname string `json:"last_name"`
	Photo   string `json:"photo_max"`
	Sex     int    `json:"sex"`
}

func NewVk(logger *zap.Logger) Vk {
	return &vkDefaultImpl{
		logger: logger,
	}
}

type vkDefaultImpl struct {
	logger *zap.Logger
}

const (
	vkURL     = "https://api.vk.com/method/users.get?"
	vkVersion = "5.131"
	timeout   = time.Second * 2
)

var ErrDoRequest = fmt.Errorf("failed to request")
var ErrValidate = fmt.Errorf("failed to validate result")

func (vk *vkDefaultImpl) GetAccount(vkToken string) (*Account, error) {
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	params := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(params)
	params.Add("access_token", vkToken)
	params.Add("fields", "photo_max,sex")
	params.Add("v", vkVersion)
	request.SetRequestURI(vkURL + params.String())
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)
	if err := fasthttp.DoTimeout(request, response, timeout); err != nil {
		vk.logger.Error("unexpected error during creating request",
			zap.Error(err),
			zap.String("access_token", vkToken),
		)
		return nil, ErrDoRequest
	}
	var vkResponse struct {
		Accounts []Account `json:"response"`
	}
	if err := json.Unmarshal(response.Body(), &vkResponse); err != nil || len(vkResponse.Accounts) == 0 {
		return nil, ErrValidate
	}
	return &vkResponse.Accounts[0], nil
}
