package websocket

import (
	"Complaingo/internal/domain/models"
	"Complaingo/internal/kafka"
	"Complaingo/internal/middleware"
	"Complaingo/internal/repository"
	"log"
	"net/http"
	"strconv"
	"sync"

	appErrors "Complaingo/internal/errors"

	"github.com/gorilla/websocket"
)

type WebsocketHandler struct {
	MessageRepo repository.MessageSaver
	kafProd     kafka.Producer
}

// client struct -- represent one connected user
type Client struct {
	Conn   *websocket.Conn
	UserID int
	Role   string
}

// channel hub struct--truck who's in what channel
type ChannelHub struct {
	subscribers map[string][]*Client
	mutex       sync.RWMutex
}

// message format of clients send and recieve
type Message struct {
	Type    string `json:"type"`
	Channel string `json:"channel"` //for pub/sub
	From    string `json:"from"`
	To      string `json:"to"` // for direct message
	Message string `json:"message"`
}

// global variables
var (
	// conver http request to ws connection
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	clients = make(map[int][]*Client) // map of connected user(userID) and their ws connections
	mutex   sync.RWMutex              // read and write mutex

	// describe who's subscribed to what channel
	hub = &ChannelHub{
		subscribers: make(map[string][]*Client),
	}
)

// add client to specific channel's subscriber list
func (h *ChannelHub) Subscribe(channel string, client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.subscribers[channel] = append(h.subscribers[channel], client)
}

// remove specific client from a single channel
func (h *ChannelHub) Unsubscribe(channel string, client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	subscribers := h.subscribers[channel]
	newClients := []*Client{}

	for _, c := range subscribers {
		if c != client {
			newClients = append(newClients, c)
		}
	}
	h.subscribers[channel] = newClients
}

func (h *ChannelHub) unsubscribeAll(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for channel, subscibers := range h.subscribers {
		newList := []*Client{}
		for _, c := range subscibers {
			if c != client {
				newList = append(newList, c)
			}
		}
		h.subscribers[channel] = newList
	}
}

// send message to all clients subscribe to a specific channel
func (h *ChannelHub) Publish(channel string, message any) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	success := true
	for _, client := range h.subscribers[channel] {
		if err := client.Conn.WriteJSON(message); err != nil {
			log.Println("Error sending to channel: ", err)
			success = false
		}
	}

	if success {
		log.Println("Message successfully sent to all subscribers.")
	}
}

// register a new connected client
func registerClient(client *Client) {
	mutex.Lock()
	defer mutex.Unlock()

	clients[client.UserID] = append(clients[client.UserID], client)
}

// unregister or disconnect
func unregisterClient(client *Client) {
	mutex.Lock()
	defer mutex.Unlock()
	clientList := clients[client.UserID]
	newClientsList := []*Client{}
	for _, c := range clientList {
		if c != client {
			newClientsList = append(newClientsList, c)
		}
	}
	clients[client.UserID] = newClientsList
}

// send message to all admins
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

// send message to user by userID
func SendToUser(userID int, message any) {
	mutex.RLock()
	defer mutex.RUnlock()

	for _, c := range clients[userID] {
		c.Conn.WriteJSON(message)
	}
}

func NewwebsocketHandler(messageRepo repository.MessageSaver, kafkaProd kafka.Producer) *WebsocketHandler {
	return &WebsocketHandler{
		MessageRepo: messageRepo,
		kafProd:     kafkaProd,
	}
}

func (h *WebsocketHandler) HandleWebsocket(w http.ResponseWriter, r *http.Request) {
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
	registerClient(client)

	// listen for messages from client
	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			log.Println("websocket read error: ", err)
			break
		}

		switch msg.Type {
		case "subscribe":
			hub.Subscribe(msg.Channel, client)
			log.Printf("User %d subscribed to %s", userID, msg.Channel)
		case "unsubscribe":
			hub.Unsubscribe(msg.Channel, client)
			log.Printf("User %d unsubscribed from %s", userID, msg.Channel)
		case "publish":
			hub.Publish(msg.Channel, msg)
			log.Printf("User %d published to %s: %s", userID, msg.Channel, msg.Message)
		case "direct":
			if msg.To == "admins" {
				go SendToAdmins(msg)

				h.kafProd.SendMessage(msg.Message)

				role := "admin"
				err := h.MessageRepo.SaveMessage(r.Context(), &models.MessageEntity{
					FromUserID: userID,
					ToUserID:   nil,
					ToRole:     &role,
					Message:    msg.Message,
				})
				if err != nil {
					middleware.WriteError(w, appErrors.ErrInvalidPayload.Wrap(err, "Failed to save"))
					return
				}

			} else {
				toID, err := strconv.Atoi(msg.To)
				if err == nil {
					go SendToUser(toID, msg)

					h.kafProd.SendMessage(msg.Message)

					err := h.MessageRepo.SaveMessage(r.Context(), &models.MessageEntity{
						FromUserID: userID,
						ToUserID:   &toID,
						ToRole:     nil,
						Message:    msg.Message,
					})
					if err != nil {
						appErrors.ErrInvalidPayload.Wrap(err, "save error")
					}
				}
				appErrors.ErrDbFailure.Wrap(err, "Invalid To field")
			}
		default:
			log.Println("unknown message type: ", msg.Type)
		}

	}

	// remove disconnected client
	hub.unsubscribeAll(client)
	unregisterClient(client)
	conn.Close()
	log.Printf("cleient %d disconnected\n", userID)
}
