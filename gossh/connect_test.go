package gossh

import (
	"bytes"
	"strings"
	"testing"
)

const (
	username = "root"
	password = "qwerty"
	ip       = "192.168.245.131"
	port     = 22
	cmd      = "cd ~/livedeploy;ls;exit"
	key      = ""
)

// const (
// 	username = "ubuntu"
// 	password = "123456qwertyAS"
// 	ip       = "101.34.216.50"
// 	port     = 22
// 	cmd      = "date"
// 	key      = ""
// )

func Test_SSH_run_simple_command(t *testing.T) {
	ciphers := []string{}
	keyExchangeList := []string{}
	session, err := connectSession(username, password, ip, key, port, ciphers, keyExchangeList)
	if err != nil {
		t.Error(err)
		return
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(cmd)
	t.Log(session.Stdout)
}

func Test_SSH_run_complex_command(t *testing.T) {
	ciphers := []string{}
	keyExchangeList := []string{}
	session, err := connectSession(username, password, ip, key, port, ciphers, keyExchangeList)
	if err != nil {
		t.Error(err)
		return
	}
	defer session.Close()

	cmdlist := strings.Split(cmd, ";")

	stdinBuf, err := session.StdinPipe()
	if err != nil {
		t.Error(err)
		return
	}

	var outbt, errbt bytes.Buffer
	session.Stdout = &outbt
	session.Stderr = &errbt

	err = session.Shell()
	if err != nil {
		t.Error(err)
		return
	}
	for _, c := range cmdlist {
		c = c + "\n"
		stdinBuf.Write([]byte(c))
		// fmt.Println()
	}

	session.Wait()
	t.Log((outbt.String() + errbt.String()))
}
