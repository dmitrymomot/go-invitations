package invitations

import (
	"errors"
	"time"

	"github.com/dmitrymomot/go-signature"
	"github.com/google/uuid"
)

type (
	// Interactor structure
	Interactor struct {
		repository InvitationRepository
	}

	// InvitationRepository interface
	InvitationRepository interface {
		GetByID(id string) (*Invitation, error)
		GetList(...Condition) ([]*Invitation, error)
		Insert(*Invitation) error
		Update(*Invitation) error
		Delete(id string) error
	}
)

// NewInteractor factory
func NewInteractor(r InvitationRepository, signingKey string) *Interactor {
	if signingKey == "" {
		signingKey = "secret%key"
	}
	signature.SetSigningKey(signingKey)
	return &Interactor{repository: r}
}

// GetByID fetch user by id
func (i *Interactor) GetByID(id string) (*Invitation, error) {
	return i.repository.GetByID(id)
}

// GetList od users with sorting and optional conditional
func (i *Interactor) GetList(c ...Condition) ([]*Invitation, error) {
	return i.repository.GetList(c...)
}

// Create new user
func (i *Interactor) Create(inv *Invitation) error {
	if inv.ID == "" {
		inv.ID = uuid.New().String()
	}
	if inv.CreatedAt == 0 {
		inv.CreatedAt = time.Now().Unix()
	}
	if inv.AccountID == "" {
		return ErrAccountIDMissed
	}
	if inv.Email == "" {
		return ErrEmailMissed
	}
	if inv.Role == "" {
		return ErrRoleMissed
	}
	if inv.InvitedBy == "" {
		return ErrInvitedByMissed
	}
	if _, err := i.repository.GetList(Email(inv.Email), AccountID(inv.AccountID)); !errors.Is(err, ErrNotFound) {
		return ErrAlreadyExists
	}
	return i.repository.Insert(inv)
}

// Update existed user
func (i *Interactor) Update(inv *Invitation) error {
	if inv.ID == "" {
		return ErrNotExistedInvitation
	}
	if inv.UpdatedAt == nil {
		t := time.Now().Unix()
		inv.UpdatedAt = &t
	}
	if inv.Role == "" {
		return ErrRoleMissed
	}
	return i.repository.Update(inv)
}

// Delete user by id
func (i *Interactor) Delete(id string) error {
	return i.repository.Delete(id)
}

// InvitationToken returns confirmation token string
func (i *Interactor) InvitationToken(inv *Invitation, ttl int64) (string, error) {
	token, err := signature.NewTemporary(inv.ID, ttl)
	if err != nil {
		return "", ErrCouldNotGenerateToken
	}
	return token, nil
}

// Accept function checks invitation token and return invitation model instance
func (i *Interactor) Accept(token string) (*Invitation, error) {
	payload, err := signature.Parse(token)
	if err != nil {
		return nil, ErrInvalidToken
	}
	id, ok := payload.(string)
	if !ok {
		return nil, ErrInvalidToken
	}
	inv, err := i.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	if inv.AcceptedAt != nil {
		return nil, ErrAlreadyAccepted
	}
	t := time.Now().Unix()
	inv.UpdatedAt = &t
	inv.AcceptedAt = &t
	if err := i.repository.Update(inv); err != nil {
		return nil, err
	}
	return inv, nil
}
