package client

import (
	"errors"
	"gosail/model"
)

var (
	usernameUnspecified = errors.New("username is not specified.")
	hostUnspecified     = errors.New("host is not specified.")
)

func checkParameterUH(host *model.SSHHost) error {
	if host.Username == "" {
		return usernameUnspecified
	}
	if host.Host == "" {
		return hostUnspecified
	}
	return nil
}
