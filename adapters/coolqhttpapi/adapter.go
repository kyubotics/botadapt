package coolqhttpapi

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/richardchien/botadapt/t"
)

type Adapter struct{}

func (adapter Adapter) UnifyEvent(event *t.Event) bool {
	// check token
	if token, ok := event.MessageSource["token"]; ok && len(token) > 0 {
		t := event.Request.Header.Get("Authorization")
		prefix := "token "
		if !strings.HasPrefix(t, prefix) || token != t[len(prefix):] {
			http.Error(event.ResponseWriter, "token is invalid", http.StatusUnauthorized)
			return false
		}
	}

	if event.Request.Header.Get("Content-Type") != "application/json" {
		http.Error(event.ResponseWriter, "invalid content type", http.StatusBadRequest)
		return false
	}

	// parse the JSON body
	jsonBytes, err := ioutil.ReadAll(event.Request.Body)
	if err != nil {
		http.Error(event.ResponseWriter, "cannot read HTTP body", http.StatusInternalServerError)
		return false
	}

	// do the unifying job in another goroutine,
	// when it's down, a JSONObject pointer will be sent to channels.EventChan
	go unifyEvent(event.MessageSource, jsonBytes)
	return true
}
