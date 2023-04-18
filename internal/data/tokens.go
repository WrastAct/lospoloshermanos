package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/WrastAct/maestro/internal/validator"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type Token struct {
	Hash   []byte    `json:"-" gorm:"not null"`
	UserID uint64    `json:"-" gorm:"not null"`
	Expiry time.Time `json:"expiry" gorm:"not null"`
	Scope  string    `json:"-" gorm:"not null"`
}

func GenerateToken(userID uint64, ttl time.Duration, scope string) (*Token, string, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, "", err
	}

	tokenPlaintext := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(tokenPlaintext))
	token.Hash = hash[:]
	return token, tokenPlaintext, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}
