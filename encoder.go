package users

import "golang.org/x/crypto/bcrypt"

// EncodeString is function to encode string with bcrypt.
func EncodeString(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CompareHash is function to compares a bcrypt hashed password with its possible plaintext equivalent.
func CompareHash(hashed, input string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(input))
}
