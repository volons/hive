package models

import (
	"sync"

	"github.com/volons/hive/libs"
)

// User model
type User struct {
	id        string `json:"id"`
	token     string `json:"-"`
	vehicleID string `json:"-"`
	name      string `json:"name"`

	permissions Permissions `json:"permissions"`

	lock *sync.RWMutex `json:"-"`
	done libs.Done     `json:"-"`
}

// NewUser creates a new user
func NewUser(id, token, vehicleID string) *User {
	return &User{
		id:          id,
		vehicleID:   vehicleID,
		token:       token,
		done:        libs.NewDone(),
		lock:        &sync.RWMutex{},
		permissions: Permissions{},
	}
}

func (u *User) Permissions() Permissions {
	return u.permissions
}

// Name returns the user's name
func (u *User) Name() string {
	u.lock.RLock()
	defer u.lock.RUnlock()

	return u.name
}

// SetName sets the user's name
func (u *User) SetName(name string) {
	u.lock.Lock()
	u.name = name
	u.lock.Unlock()
}

// ID returns the user's ID
func (u *User) ID() string {
	return u.id
}

// Token returns the user's token
func (u *User) Token() string {
	return u.token
}

// VehicleID returns the id of the vehicle to which
// the user is connected
func (u *User) VehicleID() string {
	return u.vehicleID
}

// Done returns the user's done channel
func (u *User) Done() <-chan bool {
	return u.done.WaitCh()
}

// Close closes the user's done channel
func (u *User) Close() {
	u.done.Done()
}
