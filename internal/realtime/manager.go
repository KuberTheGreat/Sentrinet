package realtime

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Manager struct {
	mu sync.Mutex
	clients map[*websocket.Conn] bool
	broadcast chan WSMessage
}

func NewManager() *Manager{
	m := &Manager{
		clients: make(map[*websocket.Conn]bool),
		broadcast: make(chan WSMessage),
	}
	go m.start()
	return m
}

func (m *Manager) Broadcast(msgType string, data interface{}){
	m.broadcast <-WSMessage{Type: msgType, Data: data}
}

func (m *Manager) start(){
	for msg := range m.broadcast{
		messageBytes, err := json.Marshal(msg)
		if err != nil{
			continue
		}

		m.mu.Lock()
		for conn := range m.clients{
			if err := conn.WriteMessage(websocket.TextMessage, messageBytes); err!= nil{
				conn.Close()
				delete(m.clients, conn)
			}
		}
		m.mu.Unlock()
	}
}

func (m *Manager) Register(c *websocket.Conn){
	m.mu.Lock()
	m.clients[c] = true
	m.mu.Unlock()
}

func (m *Manager) Unregister(c *websocket.Conn){
	m.mu.Lock()
	delete(m.clients, c)
	m.mu.Unlock()
}