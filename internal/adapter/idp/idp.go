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

