package cli

import (
	"gosail/cycle"

	"github.com/desertbit/grumble"
)

var (
	srcPath  string
	destPath string
	tar      bool
)

func init() {
	pullCommand := &grumble.Command{
		Name: "pull",
		Help: "Pull can copy file from hosts or pods, and create folders to distinguish",
		Args: func(a *grumble.Args) {
			a.String("src", "srcPath", grumble.Default(""))
			a.String("dest", "destPath", grumble.Default("."))
		},
		Flags: func(f *grumble.Flags) {
			f.Bool("", "tar", false, "tar")
			f.String("", "src", "", "srcPath")
			f.String("", "dest", ".", "destPath")
		},
		Run: func(c *grumble.Context) error {
			setPullArgs(c)
			if isK8s {
				k8sPull()
			} else {
				pull()
			}
			return nil
		},
	}
	Gosail.AddCommand(pullCommand)
}

func setPullArgs(c *grumble.Context) {
	srcPath = GetValue(c, "src", "").(string)
	destPath = GetValue(c, "dest", ".").(string)
}

func pull() {
	clientConfig, _ = cycle.GetClientConfig(keyExchanges, ciphers, cmdLine, "", hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, linuxMode)
	cycle.PullAndShow(clientConfig, &srcPath, &destPath, tar)
}
