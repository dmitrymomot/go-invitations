package invitations

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type (
	// Repository structure is implementation of UserRepository interface
	Repository struct {
		db        *sqlx.DB
		tableName string
	}
)

// NewRepository factory
func NewRepository(db *sql.DB, driverName, tableName string) *Repository {
	return &Repository{db: sqlx.NewDb(db, driverName), tableName: tableName}
}

// GetByID fetch invitation record by id
func (r *Repository) GetByID(id string) (*Invitation, error) {
	q := "SELECT * FROM %s WHERE id = ? LIMIT 1"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	inv := &Invitation{}
	if err := r.db.Get(inv, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "get invitation by id")
	}
	return inv, nil
}

// GetList fetch invitations list
func (r *Repository) GetList(c ...Condition) ([]*Invitation, error) {
	q := "SELECT * FROM %s "
	q = fmt.Sprintf(q, r.tableName)
	sq, params := conditionsToQuery(c...)
	q = q + sq
	q = r.db.Rebind(q)
	ul := make([]*Invitation, 0)
	if err := r.db.Select(&ul, q, params...); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, errors.Wrap(err, "get invitations list")
	}
	return ul, nil
}

// Insert a new invitation record
func (r *Repository) Insert(inv *Invitation) error {
	q := "INSERT INTO %s (`id`, `account_id`, `email`, `role`, `invited_by`, `created_at`) VALUES (?, ?, ?, ?, ?, ?);"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	if _, err := r.db.Exec(
		q,
		inv.ID, inv.AccountID, inv.Email,
		inv.Role, inv.InvitedBy, inv.CreatedAt,
	); err != nil {
		return errors.Wrap(err, "store invitation")
	}
	return nil
}

// Update existed invitation record
func (r *Repository) Update(inv *Invitation) error {
	q := "UPDATE %s SET `role`=?, `accepted_at`=?, `updated_at`=? WHERE id=?;"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	if _, err := r.db.Exec(q, inv.Role, inv.AcceptedAt, inv.UpdatedAt, inv.ID); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return errors.Wrap(err, "update invitation")
	}
	return nil
}

// Delete a invitation record by id
func (r *Repository) Delete(id string) error {
	q := "DELETE FROM %s WHERE id=?;"
	q = fmt.Sprintf(q, r.tableName)
	q = r.db.Rebind(q)
	if _, err := r.db.Exec(q, id); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return errors.Wrap(err, "delete invitation")
	}
	return nil
}
