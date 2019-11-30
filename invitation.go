package invitations

type (
	// Invitation model structure
	Invitation struct {
		ID         string `db:"id" json:"id"`
		AccountID  string `db:"account_id" json:"account_id"`
		Email      string `db:"email" json:"email"`
		Role       string `db:"role" json:"role"`
		InvitedBy  string `db:"invited_by" json:"invited_by"`
		CreatedAt  int64  `db:"created_at" json:"created_at"`
		UpdatedAt  *int64 `db:"updated_at" json:"updated_at,omitempty"`
		AcceptedAt *int64 `db:"accepted_at" json:"accepted_at,omitempty"`
	}
)
