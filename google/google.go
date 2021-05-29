package google

import (
	"encoding/json"
	"fmt"
	"github.com/tusa-plus/core/common"
	"io/ioutil"
	"net/http"
)

type Google struct {
	httpClientPool common.HttpClientPool
	tokenType      string
}

const (
	gglUrl = "https://www.googleapis.com/userinfo/v2/me"
)

func (ggl *Google) GetEmail(gglToken string) (string, error) {
	request, err := http.NewRequest("GET", gglUrl, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	request.Header.Add("Authorization", ggl.tokenType+" "+gglToken)
	client := ggl.httpClientPool.Get()
	defer ggl.httpClientPool.Put(client)
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
		return "", fmt.Errorf("can't unmarshal ggl response body: %v", err)
	}
	var email string
	if err := json.Unmarshal(responseJson["email"], &email); err != nil {
		return "", fmt.Errorf("can't get email from body: %v %v", string(body), err)
	}
	return email, nil
}
