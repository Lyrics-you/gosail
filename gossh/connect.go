package gossh

import (
	"bytes"
	"fmt"
	"gosail/model"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

func connect(user, password, host, key string, port int, cipherList, keyExchangeList []string) (*ssh.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		config       ssh.Config
		// session      *ssh.Session
		err error
	)
	// get authentication method
	auth = make([]ssh.AuthMethod, 0)
	if key == "" {
		auth = append(auth, ssh.Password(password))
	} else {
		pemBytes, err := ioutil.ReadFile(key)
		if err != nil {
			return nil, err
		}

		var signer ssh.Signer
		if password == "" {
			signer, err = ssh.ParsePrivateKey(pemBytes)
		} else {
			signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(password))
		}
		if err != nil {
			return nil, err
		}
		auth = append(auth, ssh.PublicKeys(signer))
	}
	if len(cipherList) == 0 {
		config.Ciphers = []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"}
	} else {
		config.Ciphers = cipherList
	}

	if len(keyExchangeList) == 0 {
		config.KeyExchanges = []string{"diffie-hellman-group-exchange-sha1", "diffie-hellman-group1-sha1", "diffie-hellman-group-exchange-sha256"}
	} else {
		config.KeyExchanges = keyExchangeList
	}

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	return client, nil
}

func connectSession(username, password, host, key string, port int, cipherList, keyExchangeList []string) (*ssh.Session, error) {
	client, err := connect(username, password, host, key, port, cipherList, keyExchangeList)
	if err != nil {
		return nil, err
	}
	// create session
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}

	setPseudoTerminal(session)

	return session, nil
}

func Dossh(username, password, host, key string, cmdlist []string, port, timeout int, cipherList, keyExchangeList []string, linuxMode bool, ch chan model.RunResult) {

	chSSH := make(chan model.RunResult)
	if linuxMode {
		go dossh_run(username, password, host, key, cmdlist, port, cipherList, keyExchangeList, chSSH)
	} else {
		go dossh_session(username, password, host, key, cmdlist, port, cipherList, keyExchangeList, chSSH)
	}
	var res model.RunResult

	select {
	case <-time.After(time.Duration(timeout) * time.Second):
		res.Host = host
		res.Username = username
		res.Success = false
		res.Result = ("ssh run timeout: " + strconv.Itoa(timeout) + " second.")
		ch <- res
	case res = <-chSSH:
		ch <- res
	}
}

func dossh_session(username, password, host, key string, cmdlist []string, port int, cipherList, keyExchangeList []string, ch chan model.RunResult) {
	session, err := connectSession(username, password, host, key, port, cipherList, keyExchangeList)

	var sshResult model.RunResult
	sshResult.Host = host
	sshResult.Username = username

	if err != nil {
		sshResult.Success = false
		sshResult.Result = fmt.Sprintf("<%s>\n", err.Error())
		ch <- sshResult
		return
	}
	defer session.Close()

	cmdlist = append(cmdlist, "exit")

	stdinBuf, _ := session.StdinPipe()

	var outbt, errbt bytes.Buffer
	session.Stdout = &outbt
	session.Stderr = &errbt

	err = session.Shell()
	if err != nil {
		sshResult.Success = false
		sshResult.Result = fmt.Sprintf("<%s>\n", err.Error())
		ch <- sshResult
		return
	}
	for _, c := range cmdlist {
		c = c + "\n"
		stdinBuf.Write([]byte(c))
	}
	session.Wait()
	if errbt.String() != "" {
		sshResult.Success = false
		sshResult.Result = errbt.String()
		ch <- sshResult
	} else {
		sshResult.Success = true
		sshResult.Result = outbt.String()
		ch <- sshResult
	}
}

func dossh_run(username, password, host, key string, cmdlist []string, port int, cipherList, keyExchangeList []string, ch chan model.RunResult) {
	session, err := connectSession(username, password, host, key, port, cipherList, keyExchangeList)
	var sshResult model.RunResult
	sshResult.Host = host
	sshResult.Username = username

	if err != nil {
		sshResult.Success = false
		sshResult.Result = fmt.Sprintf("<%s>\n", err.Error())
		ch <- sshResult
		return
	}
	defer session.Close()

	cmdlist = append(cmdlist, "exit")
	newcmd := strings.Join(cmdlist, "&&")

	var outbt, errbt bytes.Buffer
	session.Stdout = &outbt

	session.Stderr = &errbt
	err = session.Run(newcmd)
	if err != nil {
		sshResult.Success = false
		sshResult.Result = fmt.Sprintf("<%s>\n", err.Error())
		ch <- sshResult
		return
	}

	if errbt.String() != "" {
		sshResult.Success = false
		sshResult.Result = errbt.String()
		ch <- sshResult
	} else {
		sshResult.Success = true
		sshResult.Result = outbt.String()
		ch <- sshResult
	}

}

func DefaultKeyPath() string {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	keyPath := path.Join(homePath, ".ssh", "id_rsa")

	return keyPath
}
