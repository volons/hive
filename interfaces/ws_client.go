package interfaces

// WsClient interface represents a client capable of sending and receiving messages
type WsClient interface {
	SendMessage(data string) error
	Messages() chan []byte
	Done() <-chan bool
	Disconnect()
}
