package authutils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed_password), nil
}

func ComparePasswords(hashed_password, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(password))
}