package clients

import "birdnest/internal/events"

type Client struct {
	EventChan chan *events.EventMessage
}

func NewClient() *Client {
	eventChan := make(chan *events.EventMessage)

	return &Client{
		EventChan: eventChan,
	}
}
