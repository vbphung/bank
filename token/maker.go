package token

import "time"

type Maker interface {
	CreateToken(email string, expiredIn time.Duration) (*Payload, string, error)
	VerifyToken(token string) (*Payload, error)
}
