package store

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/volons/hive/libs"
	"github.com/volons/hive/libs/db"
	"github.com/volons/hive/libs/pubsub"
	"github.com/volons/hive/models"
)

type users struct {
	//sync.RWMutex
	*pubsub.Topic
	//byToken     map[string]*models.User
	//byVehicleID map[string]*models.User
}

func newUsers() *users {
	return &users{
		Topic: pubsub.NewTopic(),
		//byToken:     make(map[string]*models.User),
		//byVehicleID: make(map[string]*models.User),
	}
}

func (u *users) Get(token string) *models.User {
	user := &models.User{}
	err := db.Get(u.userKey(token), user)
	if err != nil {
		log.Println(err)
		return nil
	}

	return user

	//u.RLock()
	//defer u.RUnlock()
	//return u.byToken[token]
}

// GenerateToken renerates an auth token
// with the given userID
func (u *users) GenerateToken(vehicleID string) string {
	token := libs.RandToken(8)
	userID := libs.RandToken(4)

	user := models.NewUser(
		userID,
		token,
		vehicleID,
	)

	u.Save(user)

	return token
}

// Authenticate retreives a user from a token
func (u *users) Authenticate(token string) (*models.User, error) {
	user := Users.Get(token)
	if user == nil {
		return nil, errors.New("Invalid token")
	}

	return user, nil
}

// Save saves the user to db
func (u *users) Save(user *models.User) {
	db.SetWithTTL(u.userKey(user.Token()), user, time.Minute*10)

	//u.Lock()

	//if prev := u.byVehicleID[user.VehicleID()]; prev != nil {
	//	delete(u.byToken, prev.Token())
	//	delete(u.byVehicleID, prev.VehicleID())
	//	prev.Close()
	//}

	//u.byToken[user.Token()] = user
	//u.byVehicleID[user.VehicleID()] = user

	//u.Unlock()

	u.Publish(nil)
}

func (u *users) Delete(user *models.User) {
	db.Delete(u.userKey(user.Token()))

	//u.Lock()
	//if u.byToken[user.Token()] == user {
	//	delete(u.byToken, user.Token())
	//}
	//if u.byVehicleID[user.VehicleID()] == user {
	//	delete(u.byVehicleID, user.VehicleID())
	//}
	//user.Close()
	//u.Unlock()

	u.Publish(nil)
}

func (u *users) userKey(token string) string {
	return fmt.Sprintf("pilot:%v", token)
}

// JSON returns the list of users authorized for each vehicle in a json serializable format
func (u *users) JSON() libs.JSONObject {
	return libs.JSONObject{}

	//list := libs.JSONObject{}

	//for key, val := range u.byVehicleID {
	//	list[key] = val
	//}

	//return list
}
