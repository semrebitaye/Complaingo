package notifier

import websocket "Complaingo/internal/websockets"

type RealTimeNotifier struct{}

func (n *RealTimeNotifier) SendToAdmins(message any) {
	websocket.SendToAdmins(message)
}

func (n *RealTimeNotifier) SendToUser(userID int, message any) {
	websocket.SendToUser(userID, message)
}
