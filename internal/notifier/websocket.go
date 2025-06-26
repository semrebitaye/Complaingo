package notifier

import "crud_api/internal/websocket"

type RealTimeNotifier struct{}

func (n *RealTimeNotifier) SendToAdmins(message any) {
	websocket.SendToAdmins(message)
}

func (n *RealTimeNotifier) SendToUser(userID int, message any) {
	websocket.SendToUser(userID, message)
}
