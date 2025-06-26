package websocket

import (
	"crud_api/internal/middleware"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	UserID int
	Role   string
}

var (
	clients = make(map[int][]*Client)
	mutex   sync.RWMutex
)

// define upgrade to switch HTTP to websocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection: ", err)
		return
	}

	// extract id and role from middleware
	userID := middleware.GetUserId(r.Context())
	role := middleware.GetUserRole(r.Context())

	client := &Client{Conn: conn, UserID: userID, Role: role}

	// add client to the global list
	mutex.Lock()
	clients[userID] = append(clients[userID], client)
	mutex.Unlock()

	// listen for messages from client
	type Message struct {
		To      string `json:"to"`
		From    string `json:"from"`
		Message string `json:"message"`
	}

	for {
		var msg Message
		if err := conn.ReadJSON(msg); err != nil {
			log.Println("websocket disconnected: ", err)
			break
		}

		if msg.To == "admins" {
			go SendToAdmins(msg)
		} else {
			toUserID, err := strconv.Atoi(msg.To)
			if err == nil {
				go SendToUser(toUserID, msg)
			}
		}

	}

	// listen for messages from client
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("websocket disconnected")
			break
		}
	}

	// remove disconnected client
	mutex.Lock()
	cleintsList := clients[userID]
	for i, c := range cleintsList {
		if c == client {
			cleintsList = append(cleintsList[:i], cleintsList[i+1:]...)
			break
		}
	}
	mutex.Unlock()
}

func SendToAdmins(message any) {
	mutex.RLock()
	defer mutex.RUnlock()

	for _, list := range clients {
		for _, c := range list {
			if c.Role == "admin" {
				c.Conn.WriteJSON(message)
			}
		}
	}
}

func SendToUser(userID int, message any) {
	mutex.RLock()
	defer mutex.RUnlock()

	for _, c := range clients[userID] {
		c.Conn.WriteJSON(message)
	}
}
