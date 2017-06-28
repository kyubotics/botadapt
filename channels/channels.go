package channels

import "github.com/richardchien/botadapt/t"

var (
	EventChan chan t.JSONObject = make(chan t.JSONObject, 10)
	// APIChan  chan *structs.JSONObject
)
