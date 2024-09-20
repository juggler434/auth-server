package models

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserInterface interface {
	CheckPassword(password string) (bool, error)
}

type User struct {
	ID             uuid.UUID `json:"id"`
	UserName       string    `json:"username"`
	Email          string    `json:"email"`
	HashedPassword []byte    `json:"-"`
}

type NewUserParams struct {
	UserName string
	Email    string
	Password []byte
}

func NewUser(params *NewUserParams) (*User, error) {
	hp, err := bcrypt.GenerateFromPassword(params.Password, 10)

	u := &User{
		ID:             uuid.New(),
		UserName:       params.UserName,
		Email:          params.Email,
		HashedPassword: hp,
	}

	return u, err
}

func (u *User) CheckPassword(password []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(u.HashedPassword, password)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}

	return false, err
}
