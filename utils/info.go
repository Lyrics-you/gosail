package utils

import (
	"errors"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

var (
	// errDstHostUnspecified = errors.New("dsthost is not specified")
	ErrCantGetUsername = errors.New("cant get username")
	ErrCantGetHost     = errors.New("cant get host")
	ErrCantGetWorkDir  = errors.New("cant get work dir")
)

func PresentWorkingDir() (string, error) {
	var wd string
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
		return ".", err
	}
	wd = strings.Replace(wd, "\n", "", -1) + "/"
	return wd, err
}

func PresentHost() (string, error) {
	var host string
	if runtime.GOOS == "linux" {
		cmd := exec.Command("hostname", "-i")
		out, err := cmd.CombinedOutput()
		if err != nil {
			return "", err
		}
		host = string(out)
	} else {
		return "", ErrCantGetHost
	}
	host = strings.Replace(host, "\n", "", -1)
	return host, nil
}

func PresentUser() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	username := u.Username
	return username, nil
}

func SetDeafultUserHostPath(user, host, path *string) error {
	var err error
	if *host == "" {
		*host, err = PresentHost()
		if err != nil {
			return err
		}
	}
	if *user == "" {
		*user, err = PresentUser()
		if err != nil {
			return err
		}
	}
	if *path == "" {
		*path, err = PresentWorkingDir()
		if err != nil {
			return err
		}
	}
	return nil
}