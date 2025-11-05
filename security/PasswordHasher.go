package security

import "github.com/go-crypt/crypt/algorithm/argon2"

func HashPassword(password string) (string, error) {
	hasher, err := argon2.New(argon2.WithProfileRFC9106LowMemory())

	if err != nil {
		return "", err
	}
	digest, err := hasher.Hash(password)
	if err != nil {
		return "", err
	}

	return digest.String(), nil
}
