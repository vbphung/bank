package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	AccountID int64     `json:"account_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func CreatePayload(accountId int64, expiredIn time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{
		ID:        id,
		AccountID: accountId,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(expiredIn),
	}, nil
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return errors.New("expired token")
	}

	return nil
}
