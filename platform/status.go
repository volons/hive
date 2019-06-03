package platform

// Status represents the connection status to the API
type Status struct {
	Connected bool   `json:"connected"`
	ID        string `json:"id"`
	Token     string `json:"token"`
	Err       error  `json:"error"`
}
