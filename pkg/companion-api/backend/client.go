package backend

import (
	"github.com/mcluseau/autorizo/auth"
)

// UserData is a simple user struct with paswordhash and claims
type UserData struct {
	PasswordHash string           `json:"password"`
	ExtraClaims  auth.ExtraClaims `json:"claims"`
}

// Client is the interface for all backends clients
type Client interface {
	CreateUser(id string, user *UserData) error
	UpdateUser(id string, update func(user *UserData) error) error
	DeleteUser(id string) error
}
