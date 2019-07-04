package autopilot

import (
	"time"

	"github.com/volons/hive/libs/store"
	"github.com/volons/hive/messages"
	"github.com/volons/hive/models"
)

func (ap *Autopilot) handleVehicleMessage(msg messages.Message) {
	switch msg.Type {
	case "position":
		ap.onPositionMessage(msg)
	case "battery":
		ap.onBatteryMessage(msg)
	case "status":
		ap.onStatusMessage(msg)
	default:
		ap.forwardToUser(msg)
	}
}

func (ap *Autopilot) onStatusMessage(msg messages.Message) {
	status, ok := msg.Data.(*models.Status)
	if !ok {
		return
	}

	store.Vehicles.SetStatus(ap.vehicleID, status)

	ap.forwardToUser(msg)
}

func (ap *Autopilot) onPositionMessage(msg messages.Message) {
	pos, ok := msg.Data.(*models.Position)
	if !ok {
		return
	}

	pos.Timestamp = time.Now()
	store.Vehicles.SetPosition(ap.vehicleID, pos)

	if ap.fence != nil {
		ap.fence.checkFence(*pos)
	}

	ap.forwardToUser(msg)
}

func (ap *Autopilot) onBatteryMessage(msg messages.Message) {
	batt, ok := msg.Data.(*models.Battery)
	if !ok {
		return
	}

	store.Vehicles.SetBattery(ap.vehicleID, batt)

	ap.forwardToUser(msg)
}

func (ap *Autopilot) forwardToUser(msg messages.Message) {
	ap.user.Send(msg)
}
