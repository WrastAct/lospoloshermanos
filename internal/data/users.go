package data

import (
	"errors"
	"time"

	"github.com/WrastAct/maestro/internal/validator"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

var AnonymousUser = &User{}

type User struct {
	ID          uint64        `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"-"`
	Name        string        `json:"name" gorm:"not null"`
	Birthday    string        `json:"birthday"`
	Address     string        `json:"address"`
	Email       string        `json:"email" gorm:"unique;not null"`
	Password    []byte        `json:"-" gorm:"not null"`
	Activated   bool          `json:"activated"`
	Tokens      []Token       `json:"-"`
	Permissions []Permissions `json:"-" gorm:"many2many:user_permissions;"`
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type Password struct {
	plaintext *string
	Hash      []byte
}

func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plaintext = &plaintextPassword
	p.Hash = hash
	return nil
}

func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 20, "password", "must not be more than 20 bytes long")
}

func ValidateDate(v *validator.Validator, date string) {
	_, err := time.Parse("2006-01-02", date)
	v.Check(err == nil, "date", "must be valid")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 100, "name", "must not be more than 100 bytes long")
	ValidateEmail(v, user.Email)
	ValidateDate(v, user.Birthday)

	if user.Password == nil {
		panic("missing password hash for user")
	}
}
