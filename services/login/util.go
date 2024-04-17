package login

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (hashedPassword string, err error) {
	var bcryptHashedPassword []byte
	if bcryptHashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
		return
	}
	return string(bcryptHashedPassword), nil
}

func CheckPassword(expectedHashedPassword, actualPlainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(expectedHashedPassword), []byte(actualPlainPassword))
	return err == nil
}
