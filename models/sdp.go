package models

// SessionDescription represents an sdp message
type SessionDescription struct {
	Type string `json:"type"`
	Sdp  string `json:"sdp"`
}
