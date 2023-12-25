package util

import (
	"bytes"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	Compare(hashedPassword []byte, password []byte) error
	Hash(password []byte) ([]byte, error)
}

type bcryptPasswordHasher struct{}

func NewBcryptPasswordHasher() PasswordHasher {
	return &bcryptPasswordHasher{}
}

func (ph *bcryptPasswordHasher) Compare(hashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, password)
}

func (ph *bcryptPasswordHasher) Hash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, 12)
}

type plainTextPasswordHasher struct{}

func NewPlainTextPasswordHasher() PasswordHasher {
	return &plainTextPasswordHasher{}
}

func (ph *plainTextPasswordHasher) Compare(hashedPassword []byte, password []byte) error {
	if bytes.Equal(hashedPassword, password) {
		return nil
	} else {
		return errors.New("password mismatch")
	}
}

func (ph *plainTextPasswordHasher) Hash(password []byte) ([]byte, error) {
	return password, nil
}
