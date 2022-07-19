package cli

import (
	"fmt"
	"gosail/cycle"
	"gosail/model"
	"gosail/utils"

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
	file         string
	loginCommand *grumble.Command
	clientConfig *model.ClientConfig = &model.ClientConfig{}
)
var Gosail = grumble.New(&grumble.Config{
	Name: "gosail",
	Description: `
gosail is a free and open source batch and concurrent command execution system,
designed to execute commands on multiple servers or k8s pods and get results with speed and efficiency.
You can also copy(pull or push) files by it.`,
	HistoryFile:           "/tmp/gosail.journal",
	Prompt:                "gosail » ",
	PromptColor:           color.New(color.FgBlack, color.Bold),
	HelpHeadlineColor:     color.New(color.FgBlack),
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
		f.Int("T", "numlimit", 20, "max timeout")
		f.Int("N", "timelimit", 30, "max execute number")
	},
})

func init() {
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
	Gosail.SetPrompt(fmt.Sprintf("gosail [%s] » ", file))
}

func setInitArgs() {
	// numLimit = 20
	// timeLimit = 30
	// linuxMode = false
	hostFile = cycle.LoginHost
	username = cycle.LoginUser
	password = cycle.LoginPwd
}
