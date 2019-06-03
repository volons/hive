package models

// FenceState represents the state of the vehicle
// by comparison to the fence
type FenceState struct {
	Slowed  bool `json:"slowed"`
	Outside bool `json:"outside"`
}
