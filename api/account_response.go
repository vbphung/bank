package api

import (
	"time"

	"github.com/google/uuid"
	db "github.com/vbph/bank/db/sqlc"
)

type accountRes struct {
	ID                int64     `json:"id"`
	Email             string    `json:"email"`
	Balance           int64     `json:"balance"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type tokenRes struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

type refreshTokenRes struct {
	tokenRes
	ID uuid.UUID `json:"id"`
}

func accountResponse(acc db.Account) accountRes {
	return accountRes{
		ID:                acc.ID,
		Email:             acc.Email,
		Balance:           acc.Balance,
		PasswordChangedAt: acc.PasswordChangedAt,
		CreatedAt:         acc.CreatedAt,
	}
}
