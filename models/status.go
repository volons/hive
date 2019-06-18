package models

// Status contains the vehicles status data
type Status struct {
	Armed bool `json:"armed"`
}

// NewBattery creates and returns a new battery model
func NewStatus(armed bool) Status {
	return Status{
		armed,
	}
}
