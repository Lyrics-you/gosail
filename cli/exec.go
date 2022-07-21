package cli

import (
	"fmt"
	"gosail/client"
	"gosail/cycle"
	"gosail/gokube"
	"gosail/utils"
	"strings"

	"github.com/abiosoft/readline"
	"github.com/desertbit/grumble"
)

var (
	cmdLine   string
	highlight string
)

func init() {
	execCommand := &grumble.Command{
		Name: "exec",
		Help: "Exec can execute commands concurrently and in batches on all hosts and k8s pods, no args to exec mode",
		Args: func(a *grumble.Args) {
			a.String("cmdline", "command line", grumble.Default(""))
			a.String("highlight", "bold highlight", grumble.Default(""))
		},
		Flags: func(f *grumble.Flags) {
			f.String("e", "cmdline", "", "command line")
			f.String("b", "highlight", "", "bold highlight")
		},
		Run: func(c *grumble.Context) error {
			setExecArgs(c)
			if cmdLine == "" {
				return readCommand()
			} else {
				if isK8s {
					k8sExec()
				} else {
					if linuxMode && highlight != "" {
						cmdLine = gokube.PerlHightlight(cmdLine, highlight)
					}
					exec()
				}
			}
			return nil
		},
	}
	Gosail.AddCommand(execCommand)
}

func setExecArgs(c *grumble.Context) {
	workPath = "~"
	cmdLine = GetValue(c, "cmdline", "").(string)
	highlight = GetValue(c, "highlight", "").(string)
}

func readCommand() error {
	interConfig := &readline.Config{
		// Prompt: "> ",
		// HistorySearchFold:      true,
		// DisableAutoSaveHistory: false,
		HistoryLimit: Gosail.Config().HistoryLimit,
		AutoComplete: &NoTab{},
	}
	if isK8s {
		interConfig.HistoryFile = "/tmp/gosail_k8s.journal"
	} else {
		interConfig.HistoryFile = "/tmp/gosail_exec.journal"
	}
	rl, err := readline.NewEx(interConfig)
	if err != nil {
		return err
	}
	defer rl.Close()
	// var command string
	for {
		path := utils.GetPathLastName(workPath)
		prompt = fmt.Sprintf("gosail [%s %s] exec » ", file, path)
		rl.SetPrompt(prompt)
		command, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}
		cmdLine = strings.TrimRight(command, "\r\n")
		if cmdLine == "exit" || cmdLine == "quit" {
			return nil
		}
		if cmdLine == "clear" {
			readline.ClearScreen(rl)
			return readCommand()
		}
		if isK8s {
			interK8sExec()
		} else {
			interExec()
		}
	}
	return nil
}

func exec() {
	clientConfig, err = cycle.GetClientConfig(keyExchanges, ciphers, cmdLine, "", hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, linuxMode)
	if err != nil {
		log.Error(err)
		return
	}
	sshResults := cycle.Exec(clientConfig)
	cycle.ShowExecResult(clientConfig.SshHosts, sshResults, &jsonMode, &linuxMode)
	if selection {
		client.LoginHostByID(clientConfig.SshHosts, sshResults, "")
	}
}
