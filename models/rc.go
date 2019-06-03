package models

import (
	"errors"
	"math"
	"sync"
	"time"
)

const (
	roll     = iota
	pitch    = iota
	throttle = iota
	yaw      = iota
	gimbal   = iota
)

//var (
//	rollRange     = [2]uint16{1079, 1942} // Chan1
//	pitchRange    = [2]uint16{1080, 1943} // Chan2
//	throttleRange = [2]uint16{1083, 1943} // Chan3
//	yawRange      = [2]uint16{1086, 1942} // Chan4
//	gimbalRange   = [2]uint16{1000, 2000} // Chan6
//)

var ranges = map[int][2]uint16{
	roll:     [2]uint16{1000, 2000}, // Chan1: roll
	pitch:    [2]uint16{1000, 2000}, // Chan2: pitch
	throttle: [2]uint16{1000, 2000}, // Chan3
	yaw:      [2]uint16{1000, 2000},
	gimbal:   [2]uint16{1000, 2000},
}

//var (
//	rollRange     = [2]uint16{1000, 2000} // Chan1
//	pitchRange    = [2]uint16{1000, 2000} // Chan2
//	throttleRange = [2]uint16{1000, 2000} // Chan3
//	yawRange      = [2]uint16{1000, 2000} // Chan4
//	gimbalRange   = [2]uint16{1000, 2000} // Chan6
//)

// RcLimits stores min and max values for Rc channels
type RcLimits struct {
	RollMax     float64
	RollMin     float64
	PitchMax    float64
	PitchMin    float64
	ThrottleMax float64
	ThrottleMin float64
}

// NewRcLimits creates and initializes an RcLimits
// struct with default values
func NewRcLimits() RcLimits {
	return RcLimits{1, -1, 1, -1, 1, -1}
}

// Rc represents radio control pwm values
// to be sent to a vehicle
type Rc struct {
	lock     sync.RWMutex `json:"-"`
	throttle float64      `json:"throttle"`
	roll     float64      `json:"roll"`
	pitch    float64      `json:"pitch"`
	yaw      float64      `json:"yaw"`
	gimbal   float64      `json:"gimbal"`
	updated  time.Time    `json:"-"`
}

// NewNullRc creates and returns a new Rc struct with
// every value set to zero
func NewNullRc() *Rc {
	rc := &Rc{
		throttle: 0,
		roll:     0,
		pitch:    0,
		yaw:      0,
		gimbal:   0,
	}

	return rc
}

// NewRc creates and returns a new Rc struct
func NewRc(throttle, roll, pitch, yaw, gimbal float64) *Rc {
	rc := &Rc{
		throttle: throttle,
		roll:     roll,
		pitch:    pitch,
		yaw:      yaw,
		gimbal:   gimbal,
	}

	return rc
}

// Copy returns a copy of this struct
func (rc *Rc) Copy() *Rc {
	new := NewNullRc()
	new.Set(rc)
	return new
}

// Set copies rc data into this one
func (rc *Rc) Set(other *Rc) {
	rc.lock.Lock()
	rc.throttle = other.throttle
	rc.roll = other.roll
	rc.pitch = other.pitch
	rc.yaw = other.yaw
	rc.gimbal = other.gimbal
	rc.lock.Unlock()
}

// SetThrottle allows to set the trottle value in percent [0, 100]
func (rc *Rc) SetThrottle(val float64) error {
	if val < -1 || val > 1 {
		return errors.New("Invalid throttle value: not in range [-1, 1]")
	}

	//log.Printf("Setting throttle to %d%%\n", percent)

	rc.lock.Lock()
	rc.throttle = val
	rc.lock.Unlock()

	return nil
}

// Throttle returns the throttle value stored in this struct
func (rc *Rc) Throttle() float64 {
	rc.lock.RLock()
	defer rc.lock.RUnlock()
	return rc.throttle
}

// SetRoll allows to set the roll value in percent [0, 100]
func (rc *Rc) SetRoll(val float64) error {
	if val < -1 || val > 1 {
		return errors.New("Invalid roll value: not in range [-1, 1]")
	}

	//log.Printf("Setting roll to %d%%\n", percent)

	rc.lock.Lock()
	rc.roll = val
	rc.lock.Unlock()

	return nil
}

// Roll returns the roll value stored in this struct
func (rc *Rc) Roll() float64 {
	rc.lock.RLock()
	defer rc.lock.RUnlock()
	return rc.roll
}

// SetPitch allows to set the pitch value in percent [0, 100]
func (rc *Rc) SetPitch(val float64) error {
	if val < -1 || val > 1 {
		return errors.New("Invalid pitch value: not in range [-1, 1]")
	}

	//log.Printf("Setting pitch to %d%%\n", percent)

	rc.lock.Lock()
	rc.pitch = val
	rc.lock.Unlock()

	return nil
}

// Pitch returns the pitch value stored in this struct
func (rc *Rc) Pitch() float64 {
	rc.lock.RLock()
	defer rc.lock.RUnlock()
	return rc.pitch
}

// SetYaw allows to set the yaw value in percent [0, 100]
func (rc *Rc) SetYaw(val float64) error {
	if val < -1 || val > 1 {
		return errors.New("Invalid yaw value: not in range [-1, 1]")
	}

	//log.Printf("Setting yaw to %d%%\n", percent)

	rc.lock.Lock()
	rc.yaw = val
	rc.lock.Unlock()

	return nil
}

// Yaw returns the yaw value stored in this struct
func (rc *Rc) Yaw() float64 {
	rc.lock.RLock()
	defer rc.lock.RUnlock()
	return rc.yaw
}

// SetGimbal allows to set the gimbal value in percent [0, 100]
func (rc *Rc) SetGimbal(val float64) error {
	if val < -1 || val > 1 {
		return errors.New("Invalid gimbal value: not in range [-1, 1]")
	}

	//log.Printf("Setting gimbal to %d%%\n", percent)

	rc.lock.Lock()
	rc.gimbal = val
	rc.lock.Unlock()

	return nil
}

// Gimbal returns the gimbal value stored in this struct
func (rc *Rc) Gimbal() float64 {
	rc.lock.RLock()
	defer rc.lock.RUnlock()
	return rc.gimbal
}

// Updated updates the updated time to now
func (rc *Rc) Updated() {
	rc.lock.Lock()
	rc.updated = time.Now()
	rc.lock.Unlock()
}

// SinceLastUpdate returns the duration since last update
func (rc *Rc) SinceLastUpdate() time.Duration {
	rc.lock.RLock()
	defer rc.lock.RUnlock()
	return time.Since(rc.updated)
}

// SetDirection allows to set pitch and roll values according to an angular
// value representing a direction where 0 is forward 90 is right etc.
func (rc *Rc) SetDirection(deg float64, speed float64) {
	rad := deg * math.Pi / 180
	//log.Printf("pitch: %v%%\n", int(50+50*math.Cos(rad)*speed))
	//log.Printf("roll: %v%%\n", int(50+50*math.Sin(rad)*speed))
	rc.SetPitch(math.Cos(rad) * speed)
	rc.SetRoll(math.Sin(rad) * speed)
}

// Rotate rotates the pitch and roll axis by the given amount in deg
func (rc *Rc) Rotate(deg float64) {
	rad := deg * math.Pi / 180

	x := rc.roll
	y := rc.pitch
	len := math.Sqrt(x*x + y*y)
	angle := math.Atan2(y, x)

	//log.Printf("pitch: %v%%\n", int(50+50*math.Cos(rad)*speed))
	//log.Printf("roll: %v%%\n", int(50+50*math.Sin(rad)*speed))
	rc.SetRoll(math.Cos(angle+rad) * len)
	rc.SetPitch(math.Sin(angle+rad) * len)
}

// ApplyLimits ensures pitch, roll and
// throttle values are between limits
func (rc *Rc) ApplyLimits(limits RcLimits) bool {
	var changed bool

	if rc.Roll() > limits.RollMax {
		changed = true
		rc.SetRoll(limits.RollMax)
	}
	if rc.Roll() < limits.RollMin {
		changed = true
		rc.SetRoll(limits.RollMin)
	}

	if rc.Pitch() > limits.PitchMax {
		changed = true
		rc.SetPitch(limits.PitchMax)
	}
	if rc.Pitch() < limits.PitchMin {
		changed = true
		rc.SetPitch(limits.PitchMin)
	}

	if rc.Throttle() > limits.ThrottleMax {
		changed = true
		rc.SetThrottle(limits.ThrottleMax)
	}
	if rc.Throttle() < limits.ThrottleMin {
		changed = true
		rc.SetThrottle(limits.ThrottleMin)
	}

	return changed
}
