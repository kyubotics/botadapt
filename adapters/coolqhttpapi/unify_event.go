package coolqhttpapi

import (
	"encoding/json"

	"github.com/richardchien/botadapt/channels"
	"github.com/richardchien/botadapt/helpers"
	"github.com/richardchien/botadapt/t"
	"github.com/richardchien/botadapt/ucbi"
)

func unifyEvent(messageSource map[string]string, jsonBytes []byte) {
	eventBase := EventBase{}
	err := json.Unmarshal(jsonBytes, &eventBase)
	if err != nil {
		// not a valid JSON
		return
	}
	var unified t.JSONObject
	switch eventBase.PostType {
	case "message":
		unified = unifyMessageEvent(messageSource, jsonBytes)
	case "event": // equivalence of "notice" events of UCBI
		unified = unifyNoticeEvent(messageSource, jsonBytes)
	case "request":
		// do nothing, because UCBI does not require unifying "request" events
	}
	if unified != nil {
		channels.EventChan <- unified
	}
}

func unifyMessageEvent(messageSource map[string]string, jsonBytes []byte) t.JSONObject {
	rawEvent := MessageEvent{}
	json.Unmarshal(jsonBytes, &rawEvent)
	unified := t.JSONObject{}

	// fill in context
	context := t.JSONObject{}
	context["platform"] = PLATFORM
	context["via"] = VIA
	context["type"] = rawEvent.MessageType
	context["user_id"] = helpers.EnsureString(rawEvent.UserID)
	switch context["type"] {
	case "group":
		context["group_id"] = helpers.EnsureString(rawEvent.GroupID)
	case "discuss":
		context["discuss_id"] = helpers.EnsureString(rawEvent.DiscussID)
	}
	context["extra"] = t.JSONObject{"message_source_id": messageSource["id"]}

	// fill in data
	data := t.JSONObject{}
	data["type"] = rawEvent.MessageType
	data["sender_id"] = helpers.EnsureString(rawEvent.UserID)
	data["sender_name"], _ = getNickname(messageSource, helpers.EnsureString(rawEvent.UserID))
	data["sender"] = data["sender_name"] // we cannot get mark name in CoolQ

	switch rawEvent.MessageType {
	case "private":
		data["*sub_type"] = rawEvent.SubType
	case "group":
		data["group_id"] = helpers.EnsureString(rawEvent.GroupID)
		data["*anonymous"] = rawEvent.Anonymous
		if rawEvent.Anonymous != "" || rawEvent.UserID == 80000000 {
			// is an anonymous message
			data["sender_name"] = "匿名用户-" + rawEvent.Anonymous
			data["sender"] = data["sender_name"]
			data["sender_role"] = "unknown"
		} else {
			// is a normal message
			memberInfoBytes := apiGet(
				messageSource,
				"/get_group_member_info",
				map[string]string{
					"group_id": helpers.EnsureString(rawEvent.GroupID),
					"user_id":  helpers.EnsureString(rawEvent.UserID),
				},
			)
			if memberInfoBytes != nil {
				memberInfo := GroupMemberInfo{}
				json.Unmarshal(memberInfoBytes, &memberInfo)
				data["sender_role"] = memberInfo.Role
			} else {
				data["sender_role"] = "unknown"
			}
		}
	case "discuss":
		data["discuss_id"] = helpers.EnsureString(rawEvent.DiscussID)
		data["sender_role"] = "member"
	}

	rawMessage := rawEvent.Message
	message := []ucbi.MessageSegment{}
	lastType := ""
	for _, seg := range rawMessage {
		ucbiSeg := seg.toUCBI(messageSource)
		if ucbiSeg.Type == "text" && lastType == "text" {
			message[len(message)-1].Text += ucbiSeg.Text
		} else {
			message = append(message, ucbiSeg)
			lastType = ucbiSeg.Type
		}
	}

	if len(message) <= 0 {
		return nil
	}

	data["message"] = message

	// fill in root level
	unified["type"] = "message"
	unified["time"] = rawEvent.Time
	unified["context"] = context
	unified["data"] = data

	return unified
}

func unifyNoticeEvent(messageSource map[string]string, jsonBytes []byte) t.JSONObject {
	return t.JSONObject{}
}
