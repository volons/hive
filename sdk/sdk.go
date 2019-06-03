package sdk

import (
	"encoding/json"

	"github.com/volons/hive/libs"
	"github.com/volons/hive/libs/admin"
	"github.com/volons/hive/messages"
	"github.com/volons/hive/models"
	"github.com/volons/hive/nodes/vehicle"
	"github.com/volons/hive/platform"
)

var vehicleID = "native"

// VehicleI is a vehicle interface for communicating between native code and go gate
//type VehicleI interface {
//	SetListener(l VehicleListenerI)
//
//	//TakeOff(cb int64)
//	//Land(cb int64)
//	//RTL(cb int64)
//	//GoTo(lat, lon, relAlt float64, cb int64)
//	//Rc(roll, pitch, yaw, throttle, gimbal float64, cb int64)
//
//	//StartWebRTCSession(cb int64)
//	//SendSDP(typ, sdp string, cb int64)
//	//SendIceCandidate(sdpMLineIndex float64, sdpMid, candidate string, cb int64)
//
//	HandleMessage(msg string, cb int64)
//
//	//SendMessageToUser(message string, cb int64)
//}

// VehicleListenerI listens to VehicleI's events
//type VehicleListenerI interface {
//	//SetCaps(json string) error
//
//	//OnWebRTCSessionDescription(typ, desc string)
//	//OnWebRTCIceCandidate(sdpMLineIndex float64, sdpMid, desc string)
//
//	//OnBattery(voltage, current, percent float64)
//	//OnReport(report string)
//	//OnPosition(lat, lon, alt, hdg, relAlt, vx, vy, vz float64)
//
//	OnMessage(msg string)
//
//	//OnUserMessage(msg string)
//
//	Remove()
//}

type ChannelI interface {
	SetListener(l ChannelListenerI)
	HandleMessage(msg string)
}

type ChannelListenerI interface {
	OnMessage(msg string)
	Remove()
}

// StatusListenerI represents a platform's connection status listener
//type StatusListenerI interface {
//	Call(status StatusI)
//}

// ErrCallbackI represents a platform's connection status listener
//type ErrCallbackI interface {
//	Call(err error)
//}

//type UserCallbackI interface {
//	Call(err error, id string, name string)
//}

// StatusI represents the platform's connection status
//type StatusI interface {
//	ID() string
//	Token() string
//	Connected() bool
//	Error() error
//}

var initialized = false

// Initialize starts the storage goroutine
func Initialize() {
	if initialized {
		return
	}

	initialized = true
}

// AddVehicle connects a vehicle interface to the gate
func AddVehicle(i ChannelI, model, capsJSON string) (string, error) {
	caps := models.Caps{}
	err := json.Unmarshal([]byte(capsJSON), &caps)
	if err != nil {
		return "", err
	}

	ch := newSDKChannel(i, messages.NewParser(func(typ string) interface{} {
		switch typ {
		case "position", "goto":
			return &models.Position{}
		case "battery":
			return &models.Battery{}
		case "rc":
			return &models.Rc{}
		case "fence":
			return &models.Fence{}
		case "webrtc:sdp":
			return &models.SessionDescription{}
		case "webrtc:icecandidate":
			return &models.IceCandidate{}
		case "webrtc:start", "takeoff", "land", "rtl":
			return &struct{}{}
		case "caps":
			return &models.Caps{}
		default:
			return &libs.JSONObject{}
		}
	}))
	v := vehicle.New(ch)
	go v.Start(models.NewVehicle(vehicleID, "", model, caps))

	return vehicleID, nil
}

// AddAdmin connects an admin interface to the gate
func AddAdmin(i ChannelI) {
	ch := newSDKChannel(i, messages.NewParser(func(typ string) interface{} {
		switch typ {
		case "position", "goto":
			return &models.Position{}
		case "battery":
			return &models.Battery{}
		case "rc":
			return &models.Rc{}
		case "fence":
			return &models.Fence{}
		case "webrtc:sdp":
			return &models.SessionDescription{}
		case "webrtc:icecandidate":
			return &models.IceCandidate{}
		case "webrtc:start", "takeoff", "land", "rtl":
			return &struct{}{}
		case "caps":
			return &models.Caps{}
		default:
			return &libs.JSONObject{}
		}
	}))

	a := admin.New(ch)
	go a.Start(models.NewAdmin("0", ""), nil)
}

// Connect connects the gate to the API
func Connect(url string) {
	go platform.Platform.Run(url)
}

//// OpenLocation sends vehicle informations to the API and
//// sets the location as open
//func OpenLocation(json string, cb ErrCallbackI) {
//	go func() {
//		data, err := libs.ParseJSON([]byte(json))
//		if err != nil {
//			cb.Call(err)
//			return
//		}
//
//		fence := models.GetFence()
//		if fence == nil {
//			cb.Call(errors.New("Fence is not set"))
//			return
//		}
//
//		data["fence"] = fence.JSON()
//
//		err = platform.Platform.OpenLocation(data)
//		cb.Call(err)
//	}()
//}

//var connectDone chan bool

//// SetStatusListener sets the listeners for connection status updates
//func SetStatusListener(statusListener StatusListenerI) {
//	if connectDone != nil {
//		close(connectDone)
//		connectDone = nil
//	}
//	if statusListener == nil {
//		return
//	}
//
//	statusCh := platform.Platform.AddListener("status").(chan platform.Status)
//	done := make(chan bool)
//	connectDone = done
//
//	go func() {
//		for {
//			select {
//			case status := <-statusCh:
//				go statusListener.Call(status)
//			case <-done:
//				platform.Platform.RemoveListener("status", statusCh)
//				return
//			}
//		}
//	}()
//}

//// SetFence sets the current fence
//func SetFence(fence []byte) error {
//	return models.SetFenceJSON(fence)
//}

// EnableFence enables the fence
//func EnableFence() error {
//	ap := autopilot.Get(vehicleID)
//	return ap.EnableFence()
//}

// DisableFence enables the fence
//func DisableFence() {
//	ap := autopilot.Get(vehicleID)
//	ap.DisableFence()
//}

// OnUserConnected tells the autopilot to send
// the user initial data like the fence
//func OnUserConnected() {
//	ap := autopilot.Get(vehicleID)
//	ap.OnUserConnected()
//}

// SetControl ensures the fence is active and gives control
// of the vehicle to the user
//func SetControl(userID string) error {
//	user := store.Vehicles.GetAssignedUser(vehicleID)
//	if user.ID() != userID {
//		return errors.New("user does not exist")
//	}
//
//	err := ap.EnableFence()
//	if err != nil {
//		return err
//	}
//	err = ap.StartRcOverride()
//	if err != nil {
//		return err
//	}
//	ap.SetPilot(userID)
//
//	return nil
//}

// HasControl checks if the given user has control
// of the vehicle
//func HasControl(userID string) bool {
//	ap := autopilot.Get(vehicleID)
//	return ap.HasControl(userID)
//}

// ResetControl gives back control to the admin
//func ResetControl() {
//	ap := autopilot.Get(vehicleID)
//	ap.SetPilot("")
//	ap.StopRcOverride()
//	ap.DisableFence()
//}

//// QueueCb represents a queue callback object
//type QueueCb interface {
//	Call(listJSON string)
//}
//
//var queueDone chan bool
//
//// ListenToQueue listens to queue updates and calls the callback
//func ListenToQueue(cb QueueCb) {
//	platform.Platform.QueueSubscribe()
//	ch := platform.Platform.AddListener("queue").(chan libs.JSONArray)
//	queueDone = make(chan bool)
//	go func() {
//		for {
//			select {
//			case list := <-ch:
//				go cb.Call(list.String())
//			case <-queueDone:
//				platform.Platform.RemoveListener("queue", ch)
//				return
//			}
//
//		}
//	}()
//}
//
//// StopListeningToQueue removes the queue listener
//func StopListeningToQueue() {
//	if queueDone != nil {
//		close(queueDone)
//	}
//}
//
//// QueuePick picks a user by id in the queue and sends him through to the cockpit
//func QueuePick(userID string, cb UserCallbackI) {
//	token := store.Users.GenerateToken(vehicleID)
//
//	go func() {
//		user, err := platform.Platform.QueuePick(userID, token)
//
//		if err != nil {
//			cb.Call(err, "", "")
//			return
//		}
//
//		cb.Call(nil, user.ID(), user.Name())
//	}()
//}
//
//// QueueNext picks the first user in the queue
//func QueueNext(cb UserCallbackI) {
//	token := store.Users.GenerateToken(vehicleID)
//
//	go func() {
//		user, err := platform.Platform.QueueNext(token)
//
//		if err != nil {
//			cb.Call(err, "", "")
//			return
//		}
//
//		cb.Call(err, user.ID(), user.Name())
//	}()
//}
//
//// CallCallback calls the callback with the specified ID
//func CallCallback(id int64, err error) {
//	callCb(id, err)
//}
