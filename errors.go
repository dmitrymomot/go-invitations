package invitations

import (
	"errors"
)

// Predefined users package errors
var (
	ErrNotFound              = errors.New("user not found")
	ErrEmailMissed           = errors.New("email is missed")
	ErrRoleMissed            = errors.New("role is missed")
	ErrAccountIDMissed       = errors.New("account id is missed")
	ErrInvitedByMissed       = errors.New("invited user id is missed")
	ErrInvalidToken          = errors.New("invalid token string")
	ErrCouldNotGenerateToken = errors.New("could not generate new token")
	ErrAlreadyExists         = errors.New("invitation for this email already exists")
	ErrNotExistedInvitation  = errors.New("could not update not existed invitation")
	ErrAlreadyAccepted       = errors.New("invitation is already accepted")
)
