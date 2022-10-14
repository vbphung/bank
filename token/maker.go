package token

import "time"

type Maker interface {
	CreateToken(id string, expiredIn time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
