package ucbi

import "github.com/richardchien/botadapt/t"

type MessageSegment struct {
	Type string       `json:"type"`
	Text string       `json:"text"`
	Data t.JSONObject `json:"data"`
}
