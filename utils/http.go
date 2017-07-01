package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var client = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 5,
	},
}

func JrpcPost(ip, rpc_port, method, params string) (interface{}, error) {
	paramsStr := fmt.Sprintf(`{"id":1, "jsonrpc":"2.0", "method":"%s","params":[%s]}`, method, params)
	url := fmt.Sprintf("http://%s:%s", ip, rpc_port)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(paramsStr))
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	//fmt.Println("body :", string(body))
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	json.Unmarshal(body, &data)

	if v, ok := data["error"]; ok {
		vMap := v.(map[string]interface{})
		return nil, errors.New(fmt.Sprintf("error code:%f error message: %s", vMap["code"].(float64), vMap["message"].(string)))
	}

	return data["result"], err
}
