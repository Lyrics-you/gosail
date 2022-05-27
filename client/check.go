package client

import (
	"errors"
	"gosail/model"
)

var (
	errUsernameUnspecified = errors.New("username is not specified")
	errHostUnspecified     = errors.New("host is not specified")
)

func checkParameterUH(host *model.SSHHost) error {
	if host.Username == "" {
		return errUsernameUnspecified
	}
	if host.Host == "" {
		return errHostUnspecified
	}
	return nil
}
