package models

// Battery contains the vehicles battery data
type Battery struct {
	Voltage float64 `json:"voltage"`
	Current float64 `json:"current"`
	Percent float64 `json:"percent"`
}

// NewBattery creates and returns a new battery model
func NewBattery(voltage, current, percent float64) Battery {
	return Battery{
		voltage,
		current,
		percent,
	}
}
