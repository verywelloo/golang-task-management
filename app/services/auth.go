package service

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	//get salt

	// hashing. default coast = 10, more is more secure, but slower
	passwordWithHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwordWithHash), nil

}
