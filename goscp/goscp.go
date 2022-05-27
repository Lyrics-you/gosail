package goscp

import (
	"fmt"
	"gosail/model"
	"gosail/utils"
	"path/filepath"
	"strings"
)

func SecureCopyPull(sshHosts []model.SSHHost, surPath, dstUser, dstHost, dstPath string) error {

	err := utils.SetDeafultUserHostPath(&dstUser, &dstHost, &dstPath)
	if err != nil {
		return err
	}
	dstPath, err = GetAbsFilePath(dstPath)
	if err != nil {
		return err
	}
	for id, sshHost := range sshHosts {
		tagPath := fmt.Sprintf("%s%s/", dstPath, sshHost.Host)
		sshHosts[id].LinuxMode = true
		sshHosts[id].CmdList = []string{
			MakePullDestDirLine(dstUser, dstHost, tagPath),
			PullDestFileLine(surPath, dstUser, dstHost, tagPath),
		}
	}
	return nil
}

func SecureCopyPush(sshHosts []model.SSHHost, surUser, surHost, surPath, dstPath string) error {
	err := utils.SetDeafultUserHostPath(&surUser, &surHost, &surPath)
	if err != nil {
		return err
	}
	surPath, err = GetAbsFilePath(surPath)
	if err != nil {
		return err
	}
	for id := range sshHosts {
		sshHosts[id].LinuxMode = true
		sshHosts[id].CmdList = []string{
			MakePushDestFileLine(dstPath),
			PushDestFileLine(surUser, surHost, surPath, dstPath),
		}
	}
	return nil
}

func MakePullDestDirLine(username, dstHost, tagPath string) string {
	// host : xxx.xxx.xxx.xxx
	// tagPath : dstPath/sshHost.Host/
	// ssh username@host "mkdir -p tagPath"
	var mkdirLine string
	if strings.Contains(tagPath, dstHost) {
		mkdirLine = fmt.Sprintf("mkdir -p %s", tagPath)
	} else {
		mkdirLine = fmt.Sprintf("ssh %s@%s mkdir -p %s", username, dstHost, tagPath)
	}
	return mkdirLine
}

func PullDestFileLine(surPath, username, dstHost, tagPath string) string {
	// dstHost : xxx.xxx.xxx.xxx
	// tagPath : dstPath/sshHost.Host/
	// scp -r surPath username@dstHost:tagPath
	var copyFileLine string
	if strings.Contains(tagPath, dstHost) {
		copyFileLine = fmt.Sprintf("/bin/cp -rf %s %s && cd %s && ls", surPath, tagPath, tagPath)
	} else {
		copyFileLine = fmt.Sprintf("scp -r %s %s@%s:%s", surPath, username, dstHost, tagPath)
	}
	return copyFileLine

}

func PushDestFileLine(username, surHost, surPath, dstPth string) string {
	// dstHost : xxx.xxx.xxx.xxx
	// scp -r username@surHost:surPath dstPath
	copyFileLine := fmt.Sprintf("scp -r  %s@%s:%s %s", username, surHost, surPath, dstPth)
	return copyFileLine

}

func MakePushDestFileLine(dstPath string) string {
	// "mkdir -p tagPath"
	mkdirLine := fmt.Sprintf("mkdir -p %s", dstPath)
	return mkdirLine
}

func GetAbsFilePath(path string) (string, error) {
	abspath, err := filepath.Abs(path)
	if err != nil {
		return path, err
	}
	return abspath + "/", nil
}
