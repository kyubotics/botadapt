package coolqhttpapi

const (
	PLATFORM = "qq"
	VIA      = "coolq-http-api"
)

type EventBase struct {
	PostType string `json:"post_type"`
	Time     int32  `json:"time"`
}

type MessageEvent struct {
	EventBase
	MessageType string           `json:"message_type"`
	SubType     string           `json:"sub_type"`
	UserID      int64            `json:"user_id"`
	GroupID     int64            `json:"group_id"`
	DiscussID   int64            `json:"discuss_id"`
	Anonymous   string           `json:"anonymous"`
	Message     []MessageSegment `json:"message"`
}

type MessageSegment struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

type StrangerInfo struct {
	// currently we only need this field
	Nickname string `json:"nickname"`
}

type GroupMemberInfo struct {
	// currently we only need this field
	Role string `json:"role"`
}
