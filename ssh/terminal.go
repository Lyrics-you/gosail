package ssh

import (
	"gosail/logger"
	"os"

	"github.com/nathan-fiscaletti/consolesize-go"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

var log = logger.Logger()

func setPseudoTerminal(session *ssh.Session) error {
	// set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO: 1, // enable echoing
		// ssh.ECHOCTL:       1,
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// request pseudo terminal
	width, height := consolesize.GetConsoleSize()
	if err := session.RequestPty("xterm", height, width, modes); err != nil {
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

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	// stdin, err := session.StdinPipe()
	// if err != nil {
	// 	return err
	// }
	// stdout, err := session.StdoutPipe()
	// if err != nil {
	// 	return err
	// }
	// stderr, err := session.StderrPipe()

	// go io.Copy(os.Stderr, stderr)
	// go io.Copy(os.Stdout, stdout)
	// go func() {
	// 	buf := make([]byte, 128)
	// 	for {
	// 		n, err := os.Stdin.Read(buf)
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			return
	// 		}
	// 		if n > 0 {
	// 			_, err = stdin.Write(buf[:n])
	// 			if err != nil {
	// 				fmt.Println("Input Id: ")
	// 				return
	// 			}
	// 		}
	// 	}
	// }()

	modes := ssh.TerminalModes{
		ssh.ECHO: 1, // enable echoing
		// ssh.ECHOCTL:       0,
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Get the current terminal file descriptor for post-interaction recovery
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		log.Errorf("session.os.stdin err : %v", err)
	}
	defer terminal.Restore(fd, oldState)

	width, height := consolesize.GetConsoleSize()

	if err := session.RequestPty("xterm", height, width, modes); err != nil {
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
	os.Stdout.WriteString("Logout after pressing the key twice! ")
	os.Stdin.Read(make([]byte, 1))

	return nil
}
