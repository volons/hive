package user

import (
	"log"

	"github.com/volons/hive/libs"
	"github.com/volons/hive/libs/autopilot"
	"github.com/volons/hive/messages"
	"github.com/volons/hive/models"
)

// User represents a user connection
type User struct {
	messages.Channel

	user *models.User

	autopilot *messages.Line
}

// NewUserController creates a new User controller
func NewUserController(ch messages.Channel) *User {
	u := &User{
		Channel:   ch,
	}

	return u
}

// Start starts the client's goroutine
func (u *User) Start(user *models.User, err error) {
	if err != nil {
		u.sendError(err, "logout")
		u.Disconnect()
		return
	}

	u.autopilot = messages.NewLine("user:" + user.ID(), true)
	u.user = user

	ap := autopilot.Get(user.VehicleID())
	ap.ConnectUser(u.autopilot)

	u.Send(messages.New("update:login", libs.JSONObject{
		"id":          u.user.ID(),
		"permissions": u.user.Permissions(),
	}))

	u.run()
}

func (u *User) run() {
	for {
		select {
		case msg := <-u.autopilot.Recv():
			u.onVehicleMessage(msg)
		case msg := <-u.Recv():
			u.onMessage(msg)

		case <-u.user.Done():
			u.Disconnect()
		case <-u.Done():
			return
		}
	}
}

func (u *User) onVehicleMessage(msg messages.Message) {
	u.Send(msg)
}

func (u *User) onMessage(msg messages.Message) {
	err := u.autopilot.Send(msg)
	if err != nil {
		log.Println("cannot send to autopilot:", err)
	}
}

func (u *User) reply(msgID string, result *string, err error) {
	data := libs.JSONObject{
		"id": msgID,
	}

	if err != nil {
		data["error"] = err.Error()
	}
	if result != nil {
		data["result"] = *result
	}

	u.Send(messages.New("reply", data))
}

func (u *User) sendError(err error, action string) error {
	return u.Send(messages.New("error", libs.JSONObject{
		"message": err.Error(),
		"action":  action,
	}))
}
