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
	httpClientPool common.HttpClientPool
}

const (
	fbUrl = "https://graph.facebook.com/me?"
)

func (fb *Facebook) GetEmail(fbToken string) (string, error) {
	params := url.Values{}
	params.Add("fields", "email")
	params.Add("access_token", fbToken)
	request, err := http.NewRequest("GET", fbUrl+params.Encode(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	client := fb.httpClientPool.Get()
	defer fb.httpClientPool.Put(client)
	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("can't do request to fb: %v", err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("can't read fb response body: %v", err)
	}
	var responseJson map[string]json.RawMessage
	if err := json.Unmarshal(body, &responseJson); err != nil {
		return "", fmt.Errorf("can't unmarshal fb response body: %v", err)
	}
	var email string
	if err := json.Unmarshal(responseJson["email"], &email); err != nil {
		return "", fmt.Errorf("can't get email from body: %v %v", string(body), err)
	}
	return email, nil
}
