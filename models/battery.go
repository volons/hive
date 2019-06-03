package models

// IBattery is and interface representing a vehicles battery status
type IBattery interface {
	Voltage() float64
	Current() float64
	Percent() float64
}

// Battery contains the vehicles battery data
type Battery struct {
	voltage float64
	current float64
	percent float64
}

// NewBattery creates and returns a new battery model
func NewBattery(voltage, current, percent float64) Battery {
	return Battery{
		voltage,
		current,
		percent,
	}
}

// Voltage returns the vehicle's battery voltage value in V
func (bat Battery) Voltage() float64 {
	return bat.voltage
}

// Current returns how much current the vehicle is using in A
func (bat Battery) Current() float64 {
	return bat.current
}

// Percent returns the vehicle's battery percent value in %
func (bat Battery) Percent() float64 {
	return bat.percent
}
