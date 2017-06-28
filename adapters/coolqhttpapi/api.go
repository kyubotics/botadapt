package coolqhttpapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"
)

var (
	client = http.Client{}
)

func apiGet(messageSource map[string]string, api string, params map[string]string) []byte {
	if params != nil {
		queries := []string{}
		for k, v := range params {
			queries = append(queries, k+"="+url.QueryEscape(v))
		}
		if len(queries) > 0 {
			api += "?" + strings.Join(queries, "&")
		}
	}

	req, err := http.NewRequest("GET", strings.TrimRight(messageSource["api_url"], "/")+api, nil)
	if err != nil {
		return nil
	}
	token, ok := messageSource["token"]
	if ok && token != "" {
		req.Header.Set("Authorization", "token "+token)
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil
	}
	jsonBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	retcode := gjson.GetBytes(jsonBytes, "retcode")
	data := gjson.GetBytes(jsonBytes, "data")
	if !retcode.Exists() || retcode.Int() != 0 || !data.Exists() || data.Type == gjson.Null {
		return nil
	}
	return []byte(data.Raw)
}

func getNickname(messageSource map[string]string, userID string) (string, error) {
	userInfoBytes := apiGet(
		messageSource,
		"/get_stranger_info",
		map[string]string{"user_id": userID},
	)
	if userInfoBytes != nil {
		userInfo := StrangerInfo{}
		json.Unmarshal(userInfoBytes, &userInfo)
		return userInfo.Nickname, nil
	}
	return "", fmt.Errorf("cannot get nickname of id \"%v\"", userID)
}
