package idp

import (
	"context"
	"fmt"
)

var (
	ErrUserNotFound       = fmt.Errorf("user not found")
	ErrInvalidCredentials = fmt.Errorf("invalid credentials")
)

type IdentityProvider interface {
	Authenticate(ctx context.Context, user, pass string) (bool, error)
}

type InMemoryIdentityProvider struct {
	users map[string]string
}

func NewInMemoryIdentityProvider() *InMemoryIdentityProvider {
	i := &InMemoryIdentityProvider{users: make(map[string]string)}
	i.AddUser("admin", "admin")
	return i
}

func (i *InMemoryIdentityProvider) AddUser(username, password string) {
	i.users[username] = password
}

func (i *InMemoryIdentityProvider) Authenticate(ctx context.Context, user, pass string) (bool, error) {
	password, ok := i.users[user]
	if !ok {
		return false, ErrUserNotFound
	}

	if password != pass {
		return false, ErrInvalidCredentials
	}

	return true, nil
}
