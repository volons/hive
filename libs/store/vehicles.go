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

func (v vehicleList) SetPosition(vehicleID string, pos *models.Position) {
	err := db.Set(v.positionKey(vehicleID), pos)
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

func (v vehicleList) PositionsJSON(t time.Time) (libs.JSONObject, int) {
	return libs.JSONObject{}, 0
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

func (v vehicleList) positionKey(vehicleID string) string {
	return fmt.Sprintf("vehicle:position:%s", vehicleID)
}
