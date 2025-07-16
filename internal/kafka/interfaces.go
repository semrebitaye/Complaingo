package kafka

type Producer interface {
	SendMessage(msg string)
}
