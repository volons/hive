package controllers

import (
	"log"
	"net/http"

	"github.com/volons/hive/libs"
	"github.com/volons/hive/libs/admin"
	"github.com/volons/hive/libs/websocket"
	"github.com/volons/hive/models"
)

// Admin websocket connection listener
func Admin(wsclient *websocket.Client, r *http.Request) *websocket.Error {
	var token string
	if t, ok := r.URL.Query()["token"]; ok && len(t) > 0 && len(t[0]) > 0 {
		token = t[0]
	}

	model, err := authenticateAdmin(token, wsclient)
	if err != nil {
		log.Println("vehicle ws: invalid token")
		return websocket.NewError(err.Error(), 401)
	}

	adminConn := admin.New(wsclient)
	go adminConn.Start(model, nil)

	log.Printf("Admin '%v' connected\n", model.ID())

	return nil
}

// authenticateVehicle creates a vehicle from a token
func authenticateAdmin(token string, client *websocket.Client) (*models.Admin, error) {
	return models.NewAdmin(libs.RandToken(8), ""), nil
	//if token == "73tk91e2LtoSiPBzkeZalaUAG822Zis59E8f7oIZ738fT22RF1" {
	//}

	//return nil, errors.New("Invalid token")
}
