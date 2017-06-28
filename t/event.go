package t

import "net/http"

type Event struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	MessageSource  map[string]string
}
