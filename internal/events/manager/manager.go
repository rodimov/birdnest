package manager

import (
	"birdnest/internal/events"
	"birdnest/internal/events/clients"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
)

type Manager struct {
	Message       chan *events.EventMessage
	NewClients    chan *clients.Client
	ClosedClients chan *clients.Client
	ListClients   map[*clients.Client]bool
}

func StartEventServer(router *gin.Engine) chan *events.EventMessage {
	manager := newManager()

	g := router.Group("/api/events")
	g.GET("", headersMiddleware(), manager.getClientConnectionHandler(), handleEventSending)

	return manager.Message
}

func newManager() *Manager {
	m := &Manager{
		Message:       make(chan *events.EventMessage),
		NewClients:    make(chan *clients.Client),
		ClosedClients: make(chan *clients.Client),
		ListClients:   make(map[*clients.Client]bool),
	}

	go m.listen()

	return m
}

func headersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}

func (m *Manager) getClientConnectionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientChan := clients.NewClient()
		m.NewClients <- clientChan

		defer func() {
			m.ClosedClients <- clientChan
		}()

		c.Set("clientChan", clientChan)
		c.Next()
	}
}

func handleEventSending(context *gin.Context) {
	v, ok := context.Get("clientChan")

	if !ok {
		return
	}

	client, ok := v.(*clients.Client)

	if !ok {
		return
	}

	context.Stream(func(w io.Writer) bool {
		if event, ok := <-client.EventChan; ok {
			data, err := json.Marshal(event.Data)

			if err != nil {
				return false
			}

			context.SSEvent(event.EventType, string(data))
			return true
		}
		return false
	})
}

func (m *Manager) listen() {
	for {
		select {
		case client := <-m.NewClients:
			m.ListClients[client] = true
		case client := <-m.ClosedClients:
			delete(m.ListClients, client)
			close(client.EventChan)
		case eventMsg := <-m.Message:
			for client := range m.ListClients {
				client.EventChan <- eventMsg
			}
		}
	}
}
