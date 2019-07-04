package store

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/volons/hive/libs"
	"github.com/volons/hive/libs/db"
	"github.com/volons/hive/libs/pubsub"
	"github.com/volons/hive/models"
)

type vehicleList struct {
	*pubsub.Topic
	vehicles *sync.Map
}

func newVehicleList() vehicleList {
	return vehicleList{
		Topic:    pubsub.NewTopic(),
		vehicles: &sync.Map{},
	}
}

// Connected flags a vehicle as connected
func (v vehicleList) Connected(vehicleID string) {
	db.SetWithTTL(v.connectionKey(vehicleID), true, time.Second*5)
	//v.vehicles.Store(key(vehicle.ID), vehicle)
	v.Publish(nil)
}

// Disconnected flags a vehicle as disconnected
func (v vehicleList) Disconnected(vehicleID string) {
	log.Printf("Vehicle '%v' disconnected\n", vehicleID)
	db.Delete(v.connectionKey(vehicleID))
	//v.vehicles.Delete(key(vehicle.ID))
	v.Publish(nil)
}

func (v vehicleList) SetStatus(vehicleID string, status *models.Status) {
	err := db.Set(v.statusKey(vehicleID), *status)
	if err != nil {
		log.Println(err)
	}
}

func (v vehicleList) SetPosition(vehicleID string, pos *models.Position) {
	err := db.Set(v.positionKey(vehicleID), *pos)
	if err != nil {
		log.Println(err)
	}
}

func (v vehicleList) SetBattery(vehicleID string, batt *models.Battery) {
	err := db.Set(v.batteryKey(vehicleID), *batt)
	if err != nil {
		log.Println(err)
	}
}

// Get returns a vehicle by ID
func (v vehicleList) Get(id string) *models.Vehicle {
	var vehicle = &models.Vehicle{}
	err := db.Get(fmt.Sprintf("vehicle:%s", id), vehicle)
	if err != nil {
		return nil
	}

	return vehicle

	//val, ok := v.vehicles.Load(key(id))
	//if !ok {
	//	return nil
	//}
	//vehicle, ok := val.(*models.Vehicle)
	//if !ok {
	//	return nil
	//}

	//return vehicle
}

func (v vehicleList) Length() int {
	return 0

	//len := 0
	//v.vehicles.Range(func(key interface{}, val interface{}) bool {
	//	len++
	//	return true
	//})

	//return len
}

// JSON returns the list of vehicles in a json serializable format
func (v vehicleList) JSON() libs.JSONObject {
	return libs.JSONObject{}

	//list := libs.JSONObject{}

	//v.vehicles.Range(func(key interface{}, val interface{}) bool {
	//	v := val.(*models.Vehicle)
	//	list[key.(string)] = v
	//	return true
	//})

	//return list
}

type Telemetry struct {
	Status models.Status `json:"status"`
	Position models.Position `json:"position"`
	Battery models.Battery `json:"battery"`
	Timestamp time.Time `json:"timestamp"`
}

func NewTelemetry() Telemetry {
	return Telemetry{Timestamp: time.Now()}
}

func (v vehicleList) TelemetryJSON(t time.Time) (map[string]Telemetry) {
	out := make(map[string]Telemetry)

	getStatuses(t, out)
	getPositions(t, out)
	getBatteries(t, out)

	return out
}

func getStatuses(t time.Time, out map[string]Telemetry) {
	keys, err := db.Find(statusPrefix)
	if err != nil {
		log.Println(err)
		return
	}

	for _, key := range keys {
		var status models.Status
		err := db.Get(key, &status)
		if err != nil {
			log.Println(err)
		} else {
			id := key[len(statusPrefix):]

			val, ok := out[id];
			if !ok {
				val = NewTelemetry()
			}

			val.Status = status
			out[id] = val
		}
	}
}

func getPositions(t time.Time, out map[string]Telemetry) {
	keys, err := db.Find(positionPrefix)
	if err != nil {
		log.Println(err)
		return
	}

	for _, key := range keys {
		var pos models.Position
		err := db.Get(key, &pos)
		if err != nil {
			log.Println(err)
		} else {
			id := key[len(positionPrefix):]

			val, ok := out[id];
			if !ok {
				val = NewTelemetry()
			}

			val.Position = pos
			out[id] = val
		}
	}
}

func getBatteries(t time.Time, out map[string]Telemetry) {
	keys, err := db.Find(batteryPrefix)
	if err != nil {
		log.Println(err)
		return
	}

	for _, key := range keys {
		var batt models.Battery
		err := db.Get(key, &batt)
		if err != nil {
			log.Println(err)
		} else {
			id := key[len(batteryPrefix):]

			val, ok := out[id];
			if !ok {
				val = NewTelemetry()
			}

			val.Battery = batt
			out[id] = val
		}
	}
}

//GetIDs returns ths IDs of all vehicles in this list
func (v vehicleList) GetIDs() []string {
	return []string{}

	//ids := []string{}

	//v.vehicles.Range(func(key interface{}, val interface{}) bool {
	//	ids = append(ids, key.(string))
	//	return true
	//})

	//return ids
}

func (v vehicleList) connectionKey(vehicleID string) string {
	return fmt.Sprintf("vehicle:connected:%s", vehicleID)
}

var statusPrefix = "vehicle:status:"
func (v vehicleList) statusKey(vehicleID string) string {
	return fmt.Sprintf("%s%s", statusPrefix, vehicleID)
}

var positionPrefix = "vehicle:position:"
func (v vehicleList) positionKey(vehicleID string) string {
	return fmt.Sprintf("%s%s", positionPrefix, vehicleID)
}

var batteryPrefix = "vehicle:battery:"
func (v vehicleList) batteryKey(vehicleID string) string {
	return fmt.Sprintf("%s%s", batteryPrefix, vehicleID)
}
