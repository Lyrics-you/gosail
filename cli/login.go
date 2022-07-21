package cli

import (
	"fmt"
	"gosail/cycle"
	"gosail/model"
	"gosail/utils"

	"github.com/desertbit/grumble"
)

var (
	hostLine string
	ipLine   string
	ipFile   string
	err      error
	port     int
	sshHosts []model.SSHHost
)

func init() {
	loginCommand = &grumble.Command{
		Name:    "login",
		Help:    "Login host to do something",
		Aliases: []string{"select"},
		Args: func(a *grumble.Args) {
			a.String("hostfile", "sepcified hostfile", grumble.Default(""))
		},
		Flags: func(f *grumble.Flags) {
			f.String("", "hostfile", "", "hostfile")
			f.String("u", "username", "", "username")
			f.String("p", "password", "", "password")
			f.Int("", "port", 22, "port")
		},
		Run: func(c *grumble.Context) error {
			setLoginArgs(c)
			file = utils.GetPathLastName(hostFile)
			if file == "" {
				file = "(none)"
			}
			prompt = fmt.Sprintf("gosail [%s] » ", file)
			c.App.SetPrompt(fmt.Sprintf("gosail [%s] » ", file))
			hostList, _ = utils.GetHostList(&hostLine, &hostFile, &ipLine, &ipFile)
			sshHosts = cycle.MakeSshHosts(hostList, []string{}, &username, &password, &key, 22, &linuxMode)
			return nil
		},
	}
	k8sCommand := &grumble.Command{
		Name: "k8s",
		Help: "K8s master to do something, need login first",
		Args: func(a *grumble.Args) {

		},
		Flags: func(f *grumble.Flags) {
			f.String("n", "namespace", "", "namespace")
			f.String("a", "app", "", "app")
			f.String("c", "container", "", "container")
			f.String("", "label", "", "label")
			f.String("", "shell", "sh", "container shell")
		},
		Run: func(c *grumble.Context) error {
			setLoginArgs(c)
			setK8sArgs(c)
			file = utils.GetPathLastName(hostFile)
			if file == "" {
				file = "(none)"
			}
			prompt = fmt.Sprintf("gosail [%s %s/%s] » ", file, namespace, app)
			c.App.SetPrompt(prompt)
			hostList, err = utils.GetHostList(&hostLine, &hostFile, &ipLine, &ipFile)
			sshHosts = cycle.MakeSshHosts(hostList, []string{}, &username, &password, &key, 22, &linuxMode)
			return nil
		},
	}
	Gosail.AddCommand(loginCommand)
	loginCommand.AddCommand(k8sCommand)
}

func setLoginArgs(c *grumble.Context) {
	isK8s = false
	h := c.Flags.String("hostfile")
	if h != "" {
		hostFile = h
	}
	u := c.Flags.String("username")
	if u != "" {
		username = u
	}
	password = c.Flags.String("password")
	port = c.Flags.Int("port")
	setSetArgs(c)
}
