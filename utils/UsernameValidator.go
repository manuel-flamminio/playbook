package utils

import "net/mail"

func IsUsernameValid(username string) bool {
	_, err := mail.ParseAddress(username)
	return err == nil
}
