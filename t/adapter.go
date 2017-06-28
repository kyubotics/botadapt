package t

type Adapter interface {
	UnifyEvent(event *Event) bool
}
