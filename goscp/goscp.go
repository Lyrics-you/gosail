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
	for id := range scpConfig.SshHosts {
		scpConfig.SshHosts[id].LinuxMode = true
		scpConfig.SshHosts[id].CmdList = []string{
			MakePushDestFileLine(scpConfig.DstPath),
			fmt.Sprintf("cd %s && pwd", scpConfig.DstPath),
		}
	}
}

func PathTagHost(dstPath, host string) string {
	if dstPath != "" && !strings.HasSuffix(dstPath, "/") {
		dstPath += "/"
	}
	tagPath := fmt.Sprintf("%s%s/", dstPath, host)
	return tagPath
}

func MakeTagDirRun(dstPath, host string, ch chan model.RunResult, wg *sync.WaitGroup) {
	tagPath := PathTagHost(dstPath, host)
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
}

func SecureCopyPullMakeDir(scpConfig *model.SCPConfig) []model.RunResult {
	var wg sync.WaitGroup
	wg.Add(scpConfig.NumLimit)
	chs := make([]chan model.RunResult, len(scpConfig.SshHosts))
	for i, sshHosts := range scpConfig.SshHosts {
		chs[i] = make(chan model.RunResult, 1)
		go MakeTagDirRun(scpConfig.DstPath, sshHosts.Host, chs[i], &wg)
	}

	mkdirResults := []model.RunResult{}

	for _, ch := range chs {
		res := <-ch

		if res.Result != "" {
			mkdirResults = append(mkdirResults, res)
		}
	}
	return mkdirResults
}

func PushFileCmd(surPath, username, dstHost, dstPath string) *exec.Cmd {
	// dstHost : xxx.xxx.xxx.xxx
	// scp -r surPath username@dstHost:dstPath
	var copyFileCmd *exec.Cmd
	var iHost string
	var err error
	if runtime.GOOS == "linux" {
		iHost, err = utils.PresentHost()
		if err != nil {
			log.Error(err)
		}
	}

	if iHost == dstHost {
		copyFileCmd = exec.Command("/bin/cp", "-rf", surPath, dstPath)
	} else {
		copyFileCmd = exec.Command("scp", "-r", surPath, fmt.Sprintf("%s@%s:%s", username, dstHost, dstPath))
	}
	return copyFileCmd
}

func PullFileCmd(username, surHost, surPath, dstPath string) *exec.Cmd {
	// surHost : xxx.xxx.xxx.xxx
	// tagPath : dstPath/host/
	// scp -r username@surHost:surPath tagPath
	var copyFileCmd *exec.Cmd
	var iHost string
	var err error
	if runtime.GOOS == "linux" {
		iHost, err = utils.PresentHost()
		if err != nil {
			log.Error(err)
		}
	}
	if iHost == surHost {
		copyFileCmd = exec.Command("/bin/cp", "-rf", surPath, dstPath)
	} else {
		copyFileCmd = exec.Command("scp", "-r", fmt.Sprintf("%s@%s:%s", username, surHost, surPath), dstPath)
	}
	return copyFileCmd
}

func MakePushDestFileLine(dstPath string) string {
	// "mkdir -p tagPath"
	mkdirLine := fmt.Sprintf("mkdir -p %s", dstPath)
	return mkdirLine
}

func SecureCopyPushRun(surPath, username, dstHost, dstPath string, ch chan model.RunResult, wg *sync.WaitGroup) {
	scpCmd := PushFileCmd(surPath, username, dstHost, dstPath)
	out, err := scpCmd.CombinedOutput()
	if err != nil {
		ch <- model.RunResult{
			Success: false,
			Host:    dstHost,
			Result:  fmt.Sprintf("Push failed with %s: %s", err, string(out)),
		}
	} else {
		ch <- model.RunResult{
			Success: true,
			Host:    dstHost,
			Result:  fmt.Sprintf("%s Done!\n", scpCmd),
		}
	}
	wg.Done()
}

func SecureCopyPullRun(username, surHost, surPath, dstPath string, ch chan model.RunResult, wg *sync.WaitGroup) {
	scpCmd := PullFileCmd(username, surHost, surPath, dstPath)
	out, err := scpCmd.CombinedOutput()
	if err != nil {
		ch <- model.RunResult{
			Success: false,
			Host:    surHost,
			Result:  fmt.Sprintf("Pull failed with %s: %s", err, string(out)),
		}
	} else {
		ch <- model.RunResult{
			Success: true,
			Host:    surHost,
			Result:  fmt.Sprintf("%s Done!\n", scpCmd),
		}
	}
	wg.Done()
}

func LimitScpWithGroup(scpConfig *model.SCPConfig, runResults []model.RunResult) ([]model.RunResult, error) {
	var wg sync.WaitGroup
	wg.Add(scpConfig.NumLimit)
	chs := make([]chan model.RunResult, len(scpConfig.SshHosts))

	for i, sshHosts := range scpConfig.SshHosts {
		chs[i] = make(chan model.RunResult, 1)
		if !runResults[i].Success {
			chs[i] <- runResults[i]
		} else {
			if scpConfig.Method == "PUSH" {
				go SecureCopyPushRun(scpConfig.SurPath, sshHosts.Username, sshHosts.Host, scpConfig.DstPath, chs[i], &wg)
			} else {
				tagPath := PathTagHost(scpConfig.DstPath, sshHosts.Host)
				go SecureCopyPullRun(sshHosts.Username, sshHosts.Host, scpConfig.SurPath, tagPath, chs[i], &wg)
			}
		}
	}

	ScpResults := []model.RunResult{}

	for _, ch := range chs {
		res := <-ch
		if res.Result != "" {
			ScpResults = append(ScpResults, res)
		}

	}
	return ScpResults, nil
}
