package events

type EventMessage struct {
	EventType string      `json:"event_type"`
	Data      interface{} `json:"data"`
}
