package autopilot

import (
	"sync"
	"time"

	"github.com/volons/hive/libs"
	"github.com/volons/hive/messages"
	"github.com/volons/hive/models"
)

var autopilots sync.Map

// Autopilot handles the controls of a vehicle,
// ensures it stays in the fence when enabled
type Autopilot struct {
	//*dispatcher.Dispatcher

	vehicle *messages.Line
	user    *messages.Line
	admin   *messages.Line

	vehicleID    string           // thread safe, set once at creation
	overridingRc *libs.AtomicBool // thread safe, atomic, set at creation

	fence *fenceHandler // not thread safe, use lock
	pilot string        // not thread safe, use lock

	autoRc   *models.Rc // not thread safe, use lock
	manualRc *models.Rc // thread safe, set once at creation
	nullRc   *models.Rc // thread safe, set once at creation
	rcTicker chan bool

	//lock *sync.RWMutex
	done libs.Done
}

type key string

// Get returns the vehicle's associated autopilot
func Get(vehicleID string) *Autopilot {
	ap := &Autopilot{}
	val, loaded := autopilots.LoadOrStore(key(vehicleID), ap)
	if loaded {
		return val.(*Autopilot)
	}

	// if new init autopilot
	//ap.Dispatcher = dispatcher.Get(vehicleID)
	ap.vehicle = messages.NewLine(false)
	ap.user = messages.NewLine(false)
	ap.admin = messages.NewLine(false)
	ap.manualRc = models.NewNullRc()
	ap.nullRc = models.NewNullRc()
	ap.overridingRc = &libs.AtomicBool{}
	ap.vehicleID = vehicleID
	//ap.lock = &sync.RWMutex{}

	go ap.run()

	return ap
}

func (ap *Autopilot) run() {
	for {
		select {
		case msg := <-ap.vehicle.Recv():
			ap.handleVehicleMessage(msg)
		case msg := <-ap.user.Recv():
			ap.handleUserMessage(msg)
		case msg := <-ap.admin.Recv():
			ap.handleAdminMessage(msg)
		case <-ap.rcTicker:
			ap.vehicle.Send(messages.New("rc", ap.GetRc()))
		case <-ap.Done():
			ap.stop()
			return
		}
	}
}

func (ap *Autopilot) stop() {
	ap.vehicle.Close()
	ap.user.Close()
	ap.admin.Close()
}

// ConnectVehicle should be called by the vehicle of the same ID to listen to events
func (ap *Autopilot) ConnectVehicle(vehicle *messages.Line) {
	ap.vehicle.Connect(vehicle)
}

// ConnectUser should be called by the user that wishes to listen to this vehicle
func (ap *Autopilot) ConnectUser(user *messages.Line) {
	ap.user.Connect(user)
}

func (ap *Autopilot) Push(msg messages.Message) error {
	return ap.admin.Push(msg)
}

// AutoRc returns the rc values calculated by the fence handler
func (ap *Autopilot) AutoRc() *models.Rc {
	return ap.autoRc
}

// SetAutoRc sets the autoRc values
func (ap *Autopilot) SetAutoRc(autoRc *models.Rc) {
	ap.autoRc = autoRc
}

// GetRc returns the appropriate rc
func (ap *Autopilot) GetRc() *models.Rc {
	autoRc := ap.AutoRc()
	if autoRc != nil {
		return autoRc
	}
	if ap.manualRc.SinceLastUpdate() < time.Second*2 {
		return ap.manualRc
	}

	return ap.nullRc
}

func (ap *Autopilot) Done() <-chan bool {
	return ap.done.WaitCh()
}

// StartOrResumeSession starts a user session
//func (ap *Autopilot) StartOrResumeSession(userID string) (chan bool, error) {
//	session := StartSession(ap.vehicleID, userID)
//	ap.user.Send(messages.New("countdown", libs.JSONObject{
//		"time": session.TimeLeft(),
//	}))
//	return session.done, nil
//}

// GetPilot gets the the current pilot userID
//func (ap *Autopilot) GetPilot() string {
//	ap.lock.Lock()
//	defer ap.lock.Unlock()
//	return ap.pilot
//}

// SetPilot sets userID as the current pilot
//func (ap *Autopilot) SetPilot(userID string) {
//	ap.lock.Lock()
//	ap.pilot = userID
//	ap.lock.Unlock()
//	ap.user.Send(messages.New("pilot", libs.JSONObject{
//		"userID": userID,
//	}))
//}

// HasControl checks if the provided user has control
// over this vehicle
//func (ap *Autopilot) HasControl(userID string) bool {
//	ap.lock.RLock()
//	defer ap.lock.RUnlock()
//	return ap.pilot == userID
//}
