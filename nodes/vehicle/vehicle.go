package vehicle

import (
	"log"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/volons/hive/libs"
	"github.com/volons/hive/libs/autopilot"
	"github.com/volons/hive/libs/callback"
	"github.com/volons/hive/libs/db"
	"github.com/volons/hive/libs/store"
	"github.com/volons/hive/messages"
	"github.com/volons/hive/models"
)

// Vehicle represents a vehicle connection
type Vehicle struct {
	messages.Channel

	vehicle *models.Vehicle
	alive   *time.Ticker

	autopilot *messages.Line
}

// New creates a new Vehicle node
func New(ch messages.Channel) *Vehicle {
	v := &Vehicle{
		Channel:   ch,
	}

	return v
}

// Start initializes a vehicle connection and starts listening to events
func (v *Vehicle) Start(id string, vehicle *models.Vehicle) {
	v.autopilot = messages.NewLine("vehicle:" + id, true)
	v.vehicle = vehicle

	v.getInfo(id)

	ap := autopilot.Get(v.vehicle.ID)
	ap.ConnectVehicle(v.autopilot)

	log.Printf("Vehicle '%v' (%v) connected\n", v.vehicle.Name, v.vehicle.ID)
	store.Vehicles.Connected(v.vehicle.ID)

	v.alive = time.NewTicker(time.Second * 2)
	v.run()
}

func (v *Vehicle) getInfo(id string) {
	cb := callback.New()
	v.Send(messages.NewRequest("info", nil, cb))
	val, err := cb.Timeout(time.Minute).Wait()
	if err != nil {
		log.Println(err)
		v.Disconnect()
		return
	}

	var vehicle models.Vehicle
	err = mapstructure.Decode(val, &vehicle)
	if err != nil {
		log.Println(err)
		v.Disconnect()
		return
	}

	vehicle.ID = id
	v.vehicle = &vehicle

	err = db.Set("vehicle:"+id, vehicle)
	if err != nil {
		log.Println("Could not save new vehicle", err)
	}
}

func (v *Vehicle) run() {
	for {
		select {
		case msg := <-v.autopilot.Recv():
			v.onUserMessage(msg)
		case msg := <-v.Recv():
			v.onMessage(msg)
		case <-v.alive.C:
			store.Vehicles.Connected(v.vehicle.ID)

		case <-v.Done():
			v.alive.Stop()
			store.Vehicles.Disconnected(v.vehicle.ID)
			return
		}
	}
}

func (v *Vehicle) onUserMessage(msg messages.Message) {
	log.Printf("forwarding message to vehicle '%v'\n", v.vehicle.ID)
	v.Send(msg)
}

func (v *Vehicle) onMessage(msg messages.Message) {
	err := v.autopilot.Send(msg)
	if err != nil {
		log.Println("Vehicle could not send message to autopilot:", err)
	}
}

func (v *Vehicle) sendError(err error, action string) error {
	return v.Send(messages.New("error", libs.JSONObject{
		"message": err.Error(),
		"action":  action,
	}))
}
