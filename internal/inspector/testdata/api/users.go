package api

import (
	"context"

	"example.com/test/shared"
)

//edge:route /users/register
type RegistrationHandler struct{}

func (h *RegistrationHandler) Post(ctx context.Context, user *shared.UserProfile) error {
	return nil
}
