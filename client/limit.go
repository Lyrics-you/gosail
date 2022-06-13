package client

import (
	"fmt"
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
			log.Warnf("%s connect error, %v", host.Host, err)
			chs[i] <- model.RunResult{
				Host:    host.Host,
				Success: false,
				Result:  fmt.Sprintf("%s connect error, %v\n", host.Host, err),
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

func clinetSSHWithGroup(host model.SSHHost, clientConfig *model.ClientConfig, ch chan model.RunResult, wg *sync.WaitGroup) error {
	ssh.Dossh(host.Username, host.Password, host.Host, host.Key, host.CmdList, host.Port,
		clientConfig.TimeLimit, clientConfig.CipherList, clientConfig.KeyExchangeList, host.LinuxMode, ch)
	wg.Done()
	return nil
}

func LimitShhWithGroup(clientConfig *model.ClientConfig) ([]model.RunResult, error) {
	var wg sync.WaitGroup
	wg.Add(clientConfig.NumLimit)
	chs := make([]chan model.RunResult, len(clientConfig.SshHosts))

	for i, host := range clientConfig.SshHosts {
		chs[i] = make(chan model.RunResult, 1)

		err := checkParameterUH(&host)
		if err != nil {
			log.Warnf("%s connect error, %v", host.Host, err)
			chs[i] <- model.RunResult{
				Host:    host.Host,
				Success: false,
				Result:  fmt.Sprintf("%s connect error, %v\n", host.Host, err),
			}
		} else {
			go clinetSSHWithGroup(host, clientConfig, chs[i], &wg)
		}
	}

	sshResults := []model.RunResult{}

	for _, ch := range chs {
		res := <-ch
		if res.Result != "" {
			sshResults = append(sshResults, res)
		}

	}
	// wg.Wait()
	return sshResults, nil
}
