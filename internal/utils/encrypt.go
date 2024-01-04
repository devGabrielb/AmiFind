package utils

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(password string) ([]byte, error) {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}
	return pass, nil
}
