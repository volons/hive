package interfaces

import (
	"github.com/volons/hive/models"
)

// User model
type User interface {
	ID() string
	Token() string
	Permissions() models.Permissions
	VehicleID() string
	Name() string
	Done() <-chan bool
	Close()
}
