package websocket

import (
	"errors"
	"log"
	"net"
	"time"

	"github.com/volons/hive/libs"
	"github.com/volons/hive/libs/callback"
	"github.com/volons/hive/messages"

	"github.com/gorilla/websocket"
)

// Client represents a websocket client connection
type Client struct {
	conn      *websocket.Conn
	outgoing  chan []byte
	incoming  chan messages.Message
	callbacks callback.Map
	done      chan bool
	parser    messages.Parser

	log bool

	pingInterval time.Duration
	writeTimeout time.Duration
	readTimeout  time.Duration
}

// NewClient creates a new client
func NewClient(log bool, parser messages.Parser) *Client {
	return &Client{
		incoming:  make(chan messages.Message),
		outgoing:  make(chan []byte),
		done:      make(chan bool),
		parser:    parser,
		callbacks: callback.NewMap(),

		log: log,

		pingInterval: time.Second * 30,
		writeTimeout: time.Second,
		readTimeout:  time.Second * 55,
	}
}

// Start starts goroutines to read and write websocket messages
func (client *Client) Start(conn *websocket.Conn, isServer bool) {
	client.conn = conn
	go client.readMessages(isServer)
	go client.writeMessages(isServer)
}

func (client *Client) readMessages(isServer bool) {
	if client.Closed() {
		return
	}

	if isServer {
		client.conn.SetReadDeadline(time.Now().Add(client.readTimeout))
		client.conn.SetPongHandler(func(string) error {
			client.conn.SetReadDeadline(time.Now().Add(client.readTimeout))
			return nil
		})
	}

	for !client.Closed() {
		messageType, data, err := client.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println(err)
				client.close(err)
			}
			return
		}

		if messageType == websocket.TextMessage && len(data) > 0 {
			client.onMessage(data)
		}
	}
}

func (client *Client) onMessage(json []byte) {
	msg, err := client.parser.Parse(json)
	if err != nil {
		log.Println("error parsing json message", err)
		return
	}

	if msg.Type == "reply" {
		client.onReply(msg)
		return
	}

	if msg.IsRequest() {
		msg.Callback().Listen(func(result interface{}, err error) {
			log.Println("reply called", result, err)

			data := libs.JSONObject{
				"id": msg.ID,
			}
			if err != nil {
				data["error"] = err.Error()
			} else {
				data["result"] = result
			}

			client.Send(messages.New("reply", data))
		})
	}

	select {
	case client.incoming <- msg:
	case <-time.After(time.Second):
		log.Println("message timed out")
	}
}

func (client *Client) onReply(msg messages.Message) {
	data := msg.JSONData()
	if data == nil {
		log.Println("received reply without data")
		return
	}

	id, ok := data.GetString("id")
	if !ok {
		log.Println("received reply without id")
		return
	}

	cb := client.callbacks.Get(id)
	if cb == nil {
		log.Println("received reply for timed out or non existent request")
		return
	}

	log.Println("onReply:", data)

	if err, ok := data.GetString("error"); ok {
		cb.Reject(errors.New(err))
		return
	}

	res, _ := data.GetAny("result")
	cb.Resolve(res)
}

func (client *Client) writeMessages(isServer bool) {
	var ping <-chan time.Time

	if isServer {
		pingTicker := time.NewTicker(client.pingInterval)
		ping = pingTicker.C
		defer pingTicker.Stop()
	}

	for !client.Closed() {
		select {
		case msg := <-client.outgoing:
			client.conn.SetWriteDeadline(time.Now().Add(client.writeTimeout))
			err := client.conn.WriteMessage(websocket.TextMessage, msg)

			if err != nil {
				log.Println(err)
				client.close(err)
				return
			}
		case <-ping:
			err := client.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(client.writeTimeout))
			if err != nil {
				log.Println(err)
				client.close(err)
				return
			}
		}
	}
}

func (client *Client) close(err error) {
	if !client.Closed() {
		close(client.done)
		client.drainOutgoing()
		if client.conn != nil {
			client.conn.Close()
		}
	}
}

func (client *Client) drainOutgoing() {
	outgoing := client.outgoing
	client.outgoing = nil

	for {
		select {
		case <-outgoing:
		default:
			return
		}
	}
}

// Closed checks if the client is done
func (client *Client) Closed() bool {
	select {
	case <-client.done:
		return true
	default:
		return false
	}
}

// Disconnect disconnects the client
func (client *Client) Disconnect() {
	client.close(nil)
}

func (client *Client) Send(msg messages.Message) error {
	str, err := msg.ToJSON()
	if err != nil {
		return err
	}

	if msg.IsRequest() {
		client.callbacks.Add(msg.ID, msg.Callback())
	}

	return client.sendMessage(str)
}

// SendMessage sends a websocket message to the client
func (client *Client) sendMessage(data string) error {
	msg := []byte(data)

	select {
	case client.outgoing <- msg:
		return nil
	case <-client.done:
		return errors.New("Connection closed")
	case <-time.After(time.Millisecond * 10):
		log.Println("Message discarded not sent within 10ms")
		return errors.New("Message discarded not sent within 10ms")
	}
}

// Messages returns the incoming messages channel
func (client *Client) Recv() <-chan messages.Message {
	return client.incoming
}

// Connect starts the client by connecting
// to the specified address
func (client *Client) Connect(addr string) error {
	conn, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return err
	}

	client.Start(conn, false)
	return nil
}

// Connected returns true if this client is connected
func (client *Client) Connected() bool {
	return client.conn != nil
}

// Done returns the done channel
func (client *Client) Done() <-chan bool {
	return client.done
}

// RemoteAddr returns the ip and port of the
// remote client or server
func (client *Client) RemoteAddr() net.Addr {
	if client.conn != nil {
		return client.conn.RemoteAddr()
	}

	return nil
}
