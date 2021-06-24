package vk

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type Vk interface {
	GetAccount(vkToken string) (*Account, error)
	GetFriends(vkToken string) ([]uint64, error)
}

type Account struct {
	ID      uint64   `json:"id"`
	Name    string   `json:"first_name"`
	Surname string   `json:"last_name"`
	Photo   string   `json:"photo_max"`
	Sex     int      `json:"sex"`
	Friends []uint64 `json:"friends"`
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
	vkGetUserURL    = "https://api.vk.com/method/users.get?"
	vkGetFriendsURL = "https://api.vk.com/method/friends.get?"
	vkVersion       = "5.131"
	timeout         = time.Second * 2
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
	request.SetRequestURI(vkGetUserURL + params.String())
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

func (vk *vkDefaultImpl) GetFriends(vkToken string) ([]uint64, error) {
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)

	params := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(params)

	params.Add("access_token", vkToken)
	params.Add("v", vkVersion)
	request.SetRequestURI(vkGetFriendsURL + params.String())
	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)

	if err := fasthttp.DoTimeout(request, response, timeout); err != nil {
		vk.logger.Error("unexpected error during creating request",
			zap.Error(err),
			zap.String("access_token", vkToken),
		)
		return nil, ErrDoRequest
	}

	type Error struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
	}
	type APIResponse struct {
		Response      json.RawMessage `json:"response"`
		ResponseError Error           `json:"error"`
	}

	var vkResponse APIResponse
	if err := json.Unmarshal(response.Body(), &vkResponse); err != nil || vkResponse.ResponseError.ErrorCode != 0 {
		return nil, ErrValidate
	}

	type FriendsResponse struct {
		Count int      `json:"count"  bson:"count"`
		Items []uint64 `json:"items"  bson:"items"`
	}

	var friendsResponse FriendsResponse
	if err := json.Unmarshal(vkResponse.Response, &friendsResponse); err != nil {
		return nil, err
	}

	return friendsResponse.Items, nil
}
