package controllers

import (
	"log"
	"net/http"

	"github.com/volons/hive/libs/websocket"
	"github.com/volons/hive/models"
	"github.com/volons/hive/nodes/vehicle"
)

// Vehicle websocket connection listener
func Vehicle(wsclient *websocket.Client, r *http.Request) *websocket.Error {
	var token string
	if t, ok := r.URL.Query()["token"]; ok && len(t) > 0 && len(t[0]) > 0 {
		token = t[0]
	} else {
		log.Println("vehicle ws: no token provided")
		return websocket.NewError("Need a token", 401)
	}

	model, err := authenticateVehicle(token, wsclient)
	if err != nil {
		log.Println("vehicle ws: invalid token", err)
		return websocket.NewError(err.Error(), 401)
	}

	vehicleNode := vehicle.New(wsclient)
	go vehicleNode.Start(token, model)

	return nil
}

// authenticateVehicle creates a vehicle from a token
func authenticateVehicle(token string, client *websocket.Client) (*models.Vehicle, error) {
	var vehicle models.Vehicle

	//err := db.Get("vehicle:"+token, &vehicle)
	//if err != nil && !db.IsNotFoudError(err) {
	//	return nil, err
	//}

	return &vehicle, nil
}
