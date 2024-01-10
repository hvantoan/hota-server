package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	socket *websocket.Conn
	send   chan []byte
}

func (manager *ClientManager) start() {
	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
			}
		case message := <-manager.broadcast:
			for conn := range manager.clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}
}

func (c *Client) read(manager *ClientManager) {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			manager.unregister <- c
			c.socket.Close()
			break
		}
		manager.broadcast <- message
	}
}

func (c *Client) write() {
	defer func() {
		c.socket.Close()
	}()

	for message := range c.send {
		c.socket.WriteMessage(websocket.TextMessage, message)
	}
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (manager *ClientManager) HandleConnections(w http.ResponseWriter, r *http.Request) {
	socket, _ := upgrader.Upgrade(w, r, nil)
	client := &Client{
		socket: socket,
		send:   make(chan []byte),
	}

	manager.register <- client

	go client.read(manager)
	go client.write()
}

// Add code to handle video streaming
func (manager *ClientManager) StreamVideo(videoData []byte) {
	manager.broadcast <- videoData
}
