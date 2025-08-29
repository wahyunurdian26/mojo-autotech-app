package utils

import (
	"errors"
	"net/mail"

	"mojo-autotech/model/user_authentication"
)

func ValidateCreateAccount(req user_authentication.RegisterReq) error {

	if req.Username == "" {
		return errors.New("username is empty")
	}
	if req.Email == "" {
		return errors.New("email is empty")
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return errors.New("email is not valid")
	}
	if req.Password == "" {
		return errors.New("password is empty")
	}
	return nil
}
