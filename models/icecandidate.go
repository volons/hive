package models

// IceCandidate represents an ice candidate message
type IceCandidate struct {
	SdpMLineIndex float64 `json:"sdpMLineIndex"`
	SdpMid        string  `json:"sdpMid"`
	Candidate     string  `json:"candidate"`
}
