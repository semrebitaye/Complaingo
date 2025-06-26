package notifier

type Notifier interface {
	SendToAdmins(message any)
	SendToUser(userID int, message any)
}
