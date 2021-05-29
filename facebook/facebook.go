package facebook

import (
	"encoding/json"
	"fmt"
	"github.com/tusa-plus/core/common"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Facebook struct {
	httpClientPool *common.HttpClientPool
}

const (
	fbUrl = "https://graph.facebook.com/me?"
)

var ErrDoRequest = fmt.Errorf("failed to request")
var ErrValidate = fmt.Errorf("failed to validate result")

func (fb *Facebook) GetEmail(fbToken string) (string, error) {
	params := url.Values{}
	params.Add("fields", "email")
	params.Add("access_token", fbToken)
	request, err := http.NewRequest("GET", fbUrl+params.Encode(), nil)
	if err != nil {
		return "", ErrDoRequest
	}
	client := fb.httpClientPool.Get()
	defer fb.httpClientPool.Put(client)
	response, err := client.Do(request)
	if err != nil {
		return "", ErrDoRequest
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", ErrDoRequest
	}
	var responseJson map[string]json.RawMessage
	if err := json.Unmarshal(body, &responseJson); err != nil {
		return "", ErrValidate
	}
	var email string
	if err := json.Unmarshal(responseJson["email"], &email); err != nil {
		return "", ErrValidate
	}
	return email, nil
}
