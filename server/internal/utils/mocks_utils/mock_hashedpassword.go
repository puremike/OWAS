package mock_utils

import (
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type BcryptHashedPassword struct {
	mock.Mock
}

func (b *BcryptHashedPassword) HashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), err
}

func (b *BcryptHashedPassword) CompareHashedPassword(hashedPassword, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return err
	}
	return nil
}
