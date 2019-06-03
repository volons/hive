package models

// A Vehicle is instantiated for each connected vehicle
//type Vehicle interface {
//	ID() string
//
//	Info() libs.JSONObject
//	Caps() Caps
//
//	//TakeOff() error
//	//Land() error
//	//Rc(*Rc) error
//	//GoTo(Position) error
//	//RTL() error
//
//	// Execute a json command:
//	// { "type": "takeoff" }
//	// { "type": "land" }
//	// { "type": "rtl" }
//	// { "type": "rc", "data": { "roll": 0, "pitch": 0, "yaw": 0, "throttle": 0, "gimbal": 0 } }
//	// { "type": "goto", "data": { "lat": 0, "lon": 0, "alt": 10 } }
//	// { "type": "goto", "data": { "lat": 0, "lon": 0, "relAlt": 10 } }
//	// { "type": "webrtc:start" }
//	// { "type": "webrtc:sdp", "data": { ... } }
//	// { "type": "webrtc:icecandidate", "data": { ... } }
//	// { "type": "webrtc:forward", "data": { "msg": "{ \"type\": \"something\" }" } }
//	//Do(cmd Message) error
//
//	//SendMessageToUser(msg string) error
//
//	//SignalingChannel() interfaces.SignalingChannel
//
//	//Done() <-chan bool
//}

type Vehicle struct {
	ID    string `json:"id"`
	token string `json:"-"`
	Name  string `json:"name"`
	Model string `json:"model"`
	Caps  Caps   `json:"-"`
}

func NewVehicle(id, token, model string, caps Caps) *Vehicle {
	return &Vehicle{
		ID:    id,
		token: token,
		Model: model,
		Caps:  caps,
	}
}
