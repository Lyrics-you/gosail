package goscp

import (
	"fmt"
	"gosail/logger"
	"gosail/model"
	"gosail/utils"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var (
	log = logger.Logger()
)

func SecureCopyPushMakeDir(scpConfig *model.SCPConfig) {
	for i := range scpConfig.SshHosts {
		scpConfig.SshHosts[i].LinuxMode = true
		scpConfig.SshHosts[i].CmdList = []string{
			MakePushDestFileLine(scpConfig.DestPath[i]),
			fmt.Sprintf("cd %s && pwd", scpConfig.DestPath[i]),
		}
	}
}

func PathTagHost(destPath, host string) string {
	if destPath != "" && !strings.HasSuffix(destPath, "/") {
		destPath += "/"
	}
	tagPath := fmt.Sprintf("%s%s/", destPath, host)
	return tagPath
}

func MakeTagDirRun(chLimit chan struct{}, destPath, host string, ch chan model.RunResult, wg *sync.WaitGroup) {
	tagPath := PathTagHost(destPath, host)
	err := os.MkdirAll(tagPath, 0777)

	if err != nil {
		ch <- model.RunResult{
			Success: false,
			Result:  fmt.Sprintf("MakeDir failed with %s", err),
		}
		log.Error(err)
	} else {
		ch <- model.RunResult{
			Success: true,
			Result:  fmt.Sprintf("MakeDir %s Done!\n", tagPath),
		}
	}
	wg.Done()
	<-chLimit
}

func SecureCopyPullMakeDir(scpConfig *model.SCPConfig) []model.RunResult {
	var wg sync.WaitGroup
	chLimit := make(chan struct{}, scpConfig.NumLimit)
	chs := make([]chan model.RunResult, len(scpConfig.SshHosts))
	for i, sshHosts := range scpConfig.SshHosts {
		wg.Add(1)
		chs[i] = make(chan model.RunResult, 1)
		chLimit <- struct{}{}
		go MakeTagDirRun(chLimit, scpConfig.DestPath[i], sshHosts.Host, chs[i], &wg)
	}
	wg.Wait()

	mkdirResults := []model.RunResult{}

	for _, ch := range chs {
		res := <-ch

		if res.Result != "" {
			mkdirResults = append(mkdirResults, res)
		}
	}
	return mkdirResults
}

func PushFileCmd(srcPath, username, destHost, destPath string) *exec.Cmd {
	// destHost : xxx.xxx.xxx.xxx
	// scp -r srcPath username@destHost:destPath
	var copyFileCmd *exec.Cmd
	var iHost string
	var err error
	if runtime.GOOS == "linux" {
		iHost, err = utils.PresentHost()
		if err != nil {
			log.Error(err)
		}
	}

	if iHost == destHost {
		copyFileCmd = exec.Command("/bin/cp", "-rf", srcPath, destPath)
	} else {
		copyFileCmd = exec.Command("scp", "-r", srcPath, fmt.Sprintf("%s@%s:%s", username, destHost, destPath))
	}
	return copyFileCmd
}

func PullFileCmd(username, srcHost, srcPath, destPath string) *exec.Cmd {
	// srcHost : xxx.xxx.xxx.xxx
	// tagPath : destPath/host/
	// scp -r username@srcHost:srcPath tagPath
	var copyFileCmd *exec.Cmd
	var iHost string
	var err error
	if runtime.GOOS == "linux" {
		iHost, err = utils.PresentHost()
		if err != nil {
			log.Error(err)
		}
	}
	if iHost == srcHost {
		copyFileCmd = exec.Command("/bin/cp", "-rf", srcPath, destPath)
	} else {
		copyFileCmd = exec.Command("scp", "-r", fmt.Sprintf("%s@%s:%s", username, srcHost, srcPath), destPath)
	}
	return copyFileCmd
}

func MakePushDestFileLine(destPath string) string {
	// "mkdir -p tagPath"
	mkdirLine := fmt.Sprintf("mkdir -p %s", destPath)
	return mkdirLine
}

func SecureCopyPushRun(chLimit chan struct{}, srcPath, username, destHost, destPath string, ch chan model.RunResult, wg *sync.WaitGroup) {
	scpCmd := PushFileCmd(srcPath, username, destHost, destPath)
	out, err := scpCmd.CombinedOutput()
	if err != nil {
		ch <- model.RunResult{
			Success:  false,
			Username: username,
			Host:     destHost,
			Result:   fmt.Sprintf("Push failed with %s, %s", err, string(out)),
		}
	} else {
		ch <- model.RunResult{
			Success:  true,
			Username: username,
			Host:     destHost,
			Result:   fmt.Sprintf("%s Done!\n", scpCmd),
		}
	}
	wg.Done()
	<-chLimit
}

func SecureCopyPullRun(chLimit chan struct{}, username, srcHost, srcPath, destPath string, ch chan model.RunResult, wg *sync.WaitGroup) {
	scpCmd := PullFileCmd(username, srcHost, srcPath, destPath)
	out, err := scpCmd.CombinedOutput()
	if err != nil {
		ch <- model.RunResult{
			Success:  false,
			Username: username,
			Host:     srcHost,
			Result:   fmt.Sprintf("Pull failed with %s, %s", err, string(out)),
		}
	} else {
		ch <- model.RunResult{
			Success:  true,
			Username: username,
			Host:     srcHost,
			Result:   fmt.Sprintf("%s Done!\n", scpCmd),
		}
	}
	wg.Done()
	<-chLimit
}
