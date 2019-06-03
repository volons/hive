package messages

type Channel interface {
	Send(msg Message) error
	Recv() <-chan Message
	Done() <-chan bool
	Disconnect()
}
