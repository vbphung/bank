package api

import (
	"time"

	db "github.com/vbph/bank/db/sqlc"
)

type accountRes struct {
	ID                int64     `json:"id"`
	FullName          string    `json:"full_name"`
	Balance           int64     `json:"balance"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func accountResponse(acc db.Account) accountRes {
	return accountRes{
		ID:                acc.ID,
		FullName:          acc.FullName,
		Balance:           acc.Balance,
		PasswordChangedAt: acc.PasswordChangedAt,
		CreatedAt:         acc.CreatedAt,
	}
}
