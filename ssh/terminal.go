package ssh

import (
	"gosail/logger"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

var log = logger.Logger()

func setPseudoTerminal(session *ssh.Session) error {
	// set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// request pseudo terminal
	if err := session.RequestPty("xterm", 32, 160, modes); err != nil {
		return err
	}
	return nil
}

func GetInteractiveTerminal(username, password, host, key string, port int, cipherList, keyExchangeList []string) error {
	client, err := connect(username, password, host, key, port, cipherList, keyExchangeList)
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(fd, oldState)

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 32, 160, modes); err != nil {
		log.Errorf("session.RequestPty err: %v", err)
		return err
	}

	if err := session.Shell(); err != nil {
		log.Errorf("session.Shell err: %v", err)
		return err
	}

	if err := session.Wait(); err != nil {
		log.Errorf("session.Wait: %v", err)
		return err
	}
	return nil
}
