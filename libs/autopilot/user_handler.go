package autopilot

import (
	"errors"
	"log"
	"time"

	"github.com/volons/hive/messages"
	"github.com/volons/hive/models"
)

func (ap *Autopilot) handleUserMessage(msg messages.Message) {
	switch msg.Type {
	case "rc":
		ap.onRc(msg)
	default:
		ap.forwardToVehicle(msg)
	}
}

func (ap *Autopilot) onRc(msg messages.Message) {
	rc, ok := msg.Data.(*models.Rc)
	if ok {
		ap.SetRCValues(rc)
	}
}

func (ap *Autopilot) forwardToVehicle(msg messages.Message) {
	err := ap.vehicle.Send(msg)
	if err != nil {
		log.Println("Autopilot could not send to vehicle:", err)
	}
}

//// TakeOff tells the vehcile to takeoff
//func (ap *Autopilot) TakeOff() error {
//	err := ap.StartRcOverride()
//	if err != nil {
//		log.Println(err)
//	}
//
//	_, err = ap.vehicle.Request(messages.New("takeoff", nil)).Wait(nil)
//	if err != nil {
//		return fmt.Errorf("Could not takeoff (%v)", err)
//	}
//
//	return nil
//}
//
//// Land tells the vehcile to land
//func (ap *Autopilot) Land() error {
//	if !ap.vehicle.Connected() {
//		return fmt.Errorf("Vehicle not connected")
//	}
//
//	ap.DisableFence()
//	ap.StopRcOverride()
//
//	_, err := ap.vehicle.Request(messages.New("land", nil)).Wait(nil)
//	if err != nil {
//		return fmt.Errorf("Could not land (%v)", err)
//	}
//
//	return nil
//}
//
//// RTL tells the vehicle to return to it's launch position and land
//func (ap *Autopilot) RTL() error {
//	if !ap.vehicle.Connected() {
//		return fmt.Errorf("Vehicle not connected")
//	}
//
//	ap.StopRcOverride()
//	ap.DisableFence()
//
//	// Return to land
//	_, err := ap.vehicle.Request(messages.New("rtl", nil)).Wait(nil)
//	if err != nil {
//		return fmt.Errorf("Could not RTL (%v)", err)
//	}
//
//	return nil
//}

// SetRCValues allows to control the vehicle
func (ap *Autopilot) SetRCValues(rc *models.Rc) error {
	ap.manualRc.Set(rc)
	ap.manualRc.Updated()

	return nil
}

// StartRcOverride starts the goroutine that
// periodically sends RC Override messages
func (ap *Autopilot) StartRcOverride() error {
	if !ap.vehicle.Connected() {
		return errors.New("Vehicle not connected")
	}

	if ap.overridingRc.Get() {
		return nil
	}

	ap.overridingRc.Set(true)

	go func() {
		for ap.overridingRc.Get() {
			select {
			case <-time.After(50 * time.Millisecond):
				ap.rcTicker <- true
			case <-ap.vehicle.Done():
				ap.StopRcOverride()
				return
			}
		}

		// Reset controls by setting every value to 0
		ap.vehicle.Send(messages.New("rc", nil))
		//vehicle.Rc(nil)
	}()

	return nil
}

// StopRcOverride prevents sending radio control informations to the vehicle
func (ap *Autopilot) StopRcOverride() bool {
	return ap.overridingRc.Swap(false)
}
