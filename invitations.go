package invitations

import "database/sql"

// New function is a factory which returns invitations Interactor instance with injected users repository
// Can be used as a helper to make the code shorter
func New(db *sql.DB, driverName, tableName, signingKey string) *Interactor {
	return NewInteractor(NewRepository(db, driverName, tableName), signingKey)
}
