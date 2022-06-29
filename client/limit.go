package client

import (
	"fmt"
	"gosail/goscp"
	ssh "gosail/gossh"
	"gosail/logger"
	"gosail/model"
	"sync"
)

var log = logger.Logger()

func clinetSSHWitchChan(chLimit chan struct{}, ch chan model.RunResult, host model.SSHHost, clientConfig *model.ClientConfig) error {
	ssh.Dossh(host.Username, host.Password, host.Host, host.Key, host.CmdList, host.Port,
		clientConfig.TimeLimit, clientConfig.CipherList, clientConfig.KeyExchangeList, host.LinuxMode,
		ch)
	<-chLimit
	return nil
}

func LimitShhWithChan(clientConfig *model.ClientConfig) ([]model.RunResult, error) {
	chLimit := make(chan struct{}, clientConfig.NumLimit) //control the number of concurrent visits
	chs := make([]chan model.RunResult, len(clientConfig.SshHosts))

	for i, host := range clientConfig.SshHosts {
		chs[i] = make(chan model.RunResult, 1)

		err := checkParameterUH(&host)
		if err != nil {
			// log.Warnf("%s connect error, %v", host.Host, err)
			chs[i] <- model.RunResult{
				Username: host.Username,
				Host:     host.Host,
				Success:  false,
				Result:   fmt.Sprintf("%s connect error, %v\n", host.Host, err),
			}
		} else {
			chLimit <- struct{}{}
			go clinetSSHWitchChan(chLimit, chs[i], host, clientConfig)
		}

	}

	sshResults := []model.RunResult{}

	for _, ch := range chs {
		res := <-ch
		if res.Result != "" {
			sshResults = append(sshResults, res)
		}

	}
	return sshResults, nil
}

func clinetSSHWithGroup(chLimit chan struct{}, host model.SSHHost, clientConfig *model.ClientConfig, ch chan model.RunResult, wg *sync.WaitGroup) error {
	ssh.Dossh(host.Username, host.Password, host.Host, host.Key, host.CmdList, host.Port,
		clientConfig.TimeLimit, clientConfig.CipherList, clientConfig.KeyExchangeList, host.LinuxMode, ch)
	wg.Done()
	<-chLimit
	return nil
}

func LimitShhWithGroup(clientConfig *model.ClientConfig) ([]model.RunResult, error) {
	var wg sync.WaitGroup
	chLimit := make(chan struct{}, clientConfig.NumLimit) //control the number of concurrent visits
	chs := make([]chan model.RunResult, len(clientConfig.SshHosts))

	for i, host := range clientConfig.SshHosts {

		chs[i] = make(chan model.RunResult, 1)
		err := checkParameterUH(&host)
		if err != nil {
			// log.Warnf("%s connect error, %v", host.Host, err)
			chs[i] <- model.RunResult{
				Host:     host.Host,
				Username: host.Username,
				Success:  false,
				Result:   fmt.Sprintf("%s connect error, %v\n", host.Host, err),
			}
		} else {
			wg.Add(1)
			chLimit <- struct{}{}
			go clinetSSHWithGroup(chLimit, host, clientConfig, chs[i], &wg)
		}
	}

	wg.Wait()

	sshResults := []model.RunResult{}

	for _, ch := range chs {
		res := <-ch
		if res.Result != "" {
			sshResults = append(sshResults, res)
		}
	}
	return sshResults, nil
}

func LimitScpWithGroup(scpConfig *model.SCPConfig, runResults []model.RunResult) ([]model.RunResult, error) {
	var wg sync.WaitGroup
	chLimit := make(chan struct{}, scpConfig.NumLimit) //control the number of concurrent visits
	chs := make([]chan model.RunResult, len(scpConfig.SshHosts))

	for i, host := range scpConfig.SshHosts {
		chs[i] = make(chan model.RunResult, 1)

		err := checkParameterUH(&host)
		if err != nil {
			// log.Warnf("%s connect error, %v", host.Host, err)
			chs[i] <- model.RunResult{
				Host:     host.Host,
				Username: host.Username,
				Success:  false,
				Result:   fmt.Sprintf("%s connect error, %v\n", host.Host, err),
			}
		} else if !runResults[i].Success {
			chs[i] <- runResults[i]
		} else {
			wg.Add(1)
			if scpConfig.Method == "PUSH" {
				// PUSH
				chLimit <- struct{}{}
				go goscp.SecureCopyPushRun(chLimit, scpConfig.SrcPath[i], host.Username, host.Host, scpConfig.DestPath[i], chs[i], &wg)
			} else {
				// PULL
				tagPath := goscp.PathTagHost(scpConfig.DestPath[i], host.Host)
				chLimit <- struct{}{}
				go goscp.SecureCopyPullRun(chLimit, host.Username, host.Host, scpConfig.SrcPath[i], tagPath, chs[i], &wg)
			}
		}
	}
	wg.Wait()

	ScpResults := []model.RunResult{}

	for _, ch := range chs {
		res := <-ch
		if res.Result != "" {
			ScpResults = append(ScpResults, res)
		}

	}
	return ScpResults, nil
}
