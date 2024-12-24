package data

import "time"

type Permissions []string

// ------------------------
type TokenScope string

const (
	ActivateEmailToken TokenScope = "ActivateEmailToken"
	UserAccessToken    TokenScope = "UserAccessToken"
)

type Token struct {
	Hash   []byte     `json:"hash"`
	UserID int64      `json:"user_id"`
	Expiry time.Time  `json:"expiry"`
	Scope  TokenScope `json:"scope"`
}

// ----------------
type UserModel struct {
	ID        int64
	CreatedAt time.Time
	Name      string
	Email     string
	Password  []byte
	Activated bool
	Version   int
}
