package vk

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type Vk interface {
	GetID(vkToken string) (uint64, error)
}

func NewVk(logger *zap.Logger) Vk {
	return &vkDefaultImpl{
		logger: logger,
	}
}

func NewMockVk() Vk {
	return &vkMockImpl{}
}

type vkDefaultImpl struct {
	logger *zap.Logger
}

type vkMockImpl struct{}

const (
	vkURL     = "https://api.vk.com/method/users.get?"
	vkVersion = "5.131"
	timeout   = time.Second * 2
)

var ErrDoRequest = fmt.Errorf("failed to request")
var ErrValidate = fmt.Errorf("failed to validate result")

func (vk *vkDefaultImpl) GetID(vkToken string) (uint64, error) {
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)
	params := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(params)
	params.Add("access_token", vkToken)
	params.Add("v", vkVersion)
	request.SetRequestURI(vkURL + params.String())
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)
	if err := fasthttp.DoTimeout(request, response, timeout); err != nil {
		vk.logger.Error("unexpected error during creating request",
			zap.Error(err),
			zap.String("access_token", vkToken),
		)
		return 0, ErrDoRequest
	}
	var vkResponse struct {
		Accounts []struct {
			Id uint64 `json:"id"`
		} `json:"response"`
	}
	if err := json.Unmarshal(response.Body(), &vkResponse); err != nil || len(vkResponse.Accounts) == 0 {
		return 0, ErrValidate
	}
	return vkResponse.Accounts[0].Id, nil
}

func (vk *vkMockImpl) GetID(vkToken string) (uint64, error) {
	res, err := strconv.ParseUint(vkToken, 10, 64)
	if err != nil {
		return 0, ErrValidate
	}
	return res, nil
}
