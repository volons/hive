package websocket

import (
	"log"
	"net/http"

	"github.com/volons/hive/messages"

	"github.com/gorilla/websocket"
)

// Error represents an http error
type Error struct {
	message string
	code    int
}

// NewError create a new websocket error
func NewError(msg string, code int) *Error {
	return &Error{msg, code}
}

// ConnectionListener is the type of the function passed to SetConnectionListener
type ConnectionListener func(*Client, *http.Request) *Error

var ws = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Server handles websocket connections
type Server struct {
	connectionListener ConnectionListener
	Log                bool
	parser             messages.Parser
}

func NewServer(parser messages.Parser) *Server {
	return &Server{
		parser: parser,
	}
}

// SetConnectionListener sets the function that will be called every time
// a connection to this websocket port is attempted if the listener returns
// an error the client will not be connected
func (s *Server) SetConnectionListener(connectionListener ConnectionListener) {
	s.connectionListener = connectionListener
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	client := NewClient(s.Log, s.parser)
	err := s.connectionListener(client, r)
	if err != nil {
		http.Error(w, err.message, err.code)
	} else {
		s.connect(w, r, client)
	}
}

func (s *Server) connect(w http.ResponseWriter, r *http.Request, client *Client) {
	conn, err := ws.Upgrade(w, r, nil)
	if err != nil {
		log.Println("websocket upgrade error", err)
		return
	}

	client.Start(conn, true)
}
