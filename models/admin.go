package models

import "github.com/volons/hive/libs"

// User model
type Admin struct {
	id    string
	token string
	done  libs.Done
}

// NewUser creates a new user
func NewAdmin(id, token string) *Admin {
	return &Admin{
		id:    id,
		token: token,
		done:  libs.NewDone(),
	}
}

// ID returns the user's ID
func (u *Admin) ID() string {
	return u.id
}

// Token returns the user's token
func (u *Admin) Token() string {
	return u.token
}

// Done returns the user's done channel
func (u *Admin) Done() <-chan bool {
	return u.done.WaitCh()
}

// Close closes the user's done channel
func (u *Admin) Close() {
	u.done.Done()
}
