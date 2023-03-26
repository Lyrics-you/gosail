package cli

import (
	"fmt"
	"gosail/cycle"
	"gosail/model"
	"gosail/utils"
	"os"

	"github.com/desertbit/grumble"
	"github.com/fatih/color"
)

var (
	key          string
	ciphers      string
	keyExchanges string
	numLimit     = 20
	timeLimit    = 30
	hostFile     = cycle.LoginHost
	username     = cycle.LoginUser
	password     = cycle.LoginPwd
)
var (
	file          = "(none)"
	prompt        = "gosail » "
	promptColor   = color.New(color.FgWhite, color.Bold)
	whoami        string
	loginCommand  *grumble.Command
	spinnerConfig *model.SpinConfig   = &model.SpinConfig{}
	clientConfig  *model.ClientConfig = &model.ClientConfig{}
	gosailConfig  *model.GosailConfig = &model.GosailConfig{}
)
var Gosail = grumble.New(&grumble.Config{
	Name: "gosail",
	Description: `
gosail is a free and open source batch and concurrent command execution system,
designed to execute commands on multiple servers or k8s pods and get results with speed and efficiency.
You can also copy(pull or push) files by it.`,
	// HistoryFile:           "/tmp/gosail.journal",
	Prompt:                prompt,
	PromptColor:           promptColor,
	HelpHeadlineColor:     promptColor,
	HelpHeadlineUnderline: true,
	HelpSubCommands:       true,
	Flags: func(f *grumble.Flags) {
		// login
		f.String("", "hostfile", "", "hostfile")
		f.String("u", "username", "", "username")
		f.String("p", "password", "", "password")
		// set
		f.String("K", "key", "", "id_rsa.pub key filepath")
		f.String("", "ciphers", "", "ssh ciphers")
		f.String("", "keyexchanges", "", "ssh keyexchangesx")
		f.Int("T", "timelimit", 30, "max timeout")
		f.Int("N", "numlimit", 20, "max execute number")
	},
})

func init() {
	whoami, _ = utils.PresentUser()
	whoami = utils.GetPathLastName(whoami)
	Gosail.Config().HistoryFile = fmt.Sprintf("/tmp/%s_gosail.journal", whoami)
	version := model.Historys[len(model.Historys)-1]
	Gosail.SetPrintASCIILogo(func(a *grumble.App) {
		a.Println(model.LOGO)
		a.Println(model.EMOJI["gosail"], version.Version)
	})
	setInitArgs()
	file = utils.GetPathLastName(hostFile)
	if file == "" {
		file = "(none)"
	}
	prompt = fmt.Sprintf("gosail [%s] » ", file)
	Gosail.SetPrompt(prompt)
	Gosail.SetInterruptHandler(func(a *grumble.App, count int) {
		if count >= 2 {
			a.Println("👌Finshed!")
			os.Exit(0)
		}
		a.Println("input Ctrl-c once more to exit")
	})
}

func setInitArgs() {
	// numLimit = 20
	// timeLimit = 30
	// linuxMode = false
	hostFile = cycle.LoginHost
	username = cycle.LoginUser
	password = cycle.LoginPwd
	gosailConfig = cycle.GetGosailConfiguration()
	// gConfig.PrintGosailConfiguration()
	if gosailConfig != nil {
		// client
		// keyExchanges
		// ciphers
		if gosailConfig.Client.NumLimit != 0 {
			numLimit = gosailConfig.Client.NumLimit
		}
		if gosailConfig.Client.TimeLimit != 0 {
			timeLimit = gosailConfig.Client.TimeLimit
		}
		// spinner
		spinnerConfig = gosailConfig.Spin
		// spinnerConfig.TimeOut = timeLimit
		// login
		if gosailConfig.Login.HostFile != "" {
			hostFile = gosailConfig.Login.HostFile
		}
		if gosailConfig.Login.IpFile != "" {
			ipFile = gosailConfig.Login.IpFile
		}
		if gosailConfig.Login.Username != "" {
			username = gosailConfig.Login.Username
		}
		if gosailConfig.Login.Password != "" {
			password = gosailConfig.Login.Password
		}
		if gosailConfig.Login.Port != 0 {
			port = gosailConfig.Login.Port
		} else {
			port = 22
		}

	}
}
