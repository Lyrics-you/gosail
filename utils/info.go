package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	ErrCantGetUsername = errors.New("cant get username")
	ErrCantGetHost     = errors.New("cant get host")
	ErrCantGetWorkDir  = errors.New("cant get work dir")
	ErrCmdListEmpty    = errors.New("commands is empty")
)

func GetAbsFilePath(path string) (string, error) {
	abspath, err := filepath.Abs(path)
	if err != nil {
		return path, err
	}
	return abspath + "/", nil
}

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
		// return "192.168.245.1", nil
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

func GetPathLastName(path string) string {
	sep := "/"
	if strings.Contains(path, "\\") {
		sep = "\\"
	}
	names := strings.Split(path, sep)
	i := len(names)
	for i > 0 && names[i-1] == "" {
		i -= 1
	}
	if i > 0 {
		return names[i-1]
	}
	return ""
}

func GetHostList(hostLine, hostFile, ipLine, ipFile *string) ([]string, error) {
	var hostList []string
	var err error
	if *ipFile != "" {
		hostList, err = GetIpListFromFile(*ipFile)
		if err != nil {
			// log.Errorf("load iplist error, %v", err)
			return []string{}, fmt.Errorf("load iplist error, %v", err)
		}
		return hostList, nil
	}

	if *hostFile != "" {
		hostList, err = GetString(*hostFile)
		if err != nil {
			// log.Errorf("load hostfile error, %v", err)
			return []string{}, fmt.Errorf("load hostfile error, %v", err)
		}
		return hostList, nil
	}

	if *ipLine != "" {
		hostList, err = GetIpListFromString(*ipLine)
		if err != nil {
			// log.Errorf("load iplist error, %v", err)
			return []string{}, fmt.Errorf("load hostfile error, %v", err)
		}
		return hostList, nil
	}

	if *hostLine != "" {
		hostList = SplitString(*hostLine)
	}
	return hostList, nil
}

func GetCmdList(cmdLine, cmdFile *string) ([]string, error) {
	var cmdList []string
	if *cmdFile != "" {
		cmdList, err := GetString(*cmdFile)
		if err != nil {
			// log.Errorf("load cmdfile error, %v", err)
			return []string{}, fmt.Errorf("load cmdfile error, %v", err)
		}
		if len(cmdList) == 0 {
			return []string{}, ErrCmdListEmpty
		}
		return cmdList, nil
	}

	if *cmdLine != "" {
		cmdList = SplitString(*cmdLine)
		if len(cmdList) == 0 {
			return []string{}, ErrCmdListEmpty
		}
		return cmdList, nil
	}

	return []string{}, ErrCmdListEmpty
}
