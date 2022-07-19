package cli

import (
	"gosail/cycle"

	"github.com/desertbit/grumble"
)

func init() {
	pushCommand := &grumble.Command{
		Name: "push",
		Help: "Pull can copy file to hosts concurrently, and create folders that do not exist",
		Args: func(a *grumble.Args) {
			a.String("src", "srcPath", grumble.Default(""))
			a.String("dest", "destPath", grumble.Default("."))
		},
		Flags: func(f *grumble.Flags) {
			f.String("", "src", "", "srcPath")
			f.String("", "dest", ".", "destPath")
		},
		Run: func(c *grumble.Context) error {
			setPushArgs(c)
			push()
			return nil
		},
	}
	Gosail.AddCommand(pushCommand)
}

func setPushArgs(c *grumble.Context) {
	srcPath = GetValue(c, "src", "").(string)
	destPath = GetValue(c, "dest", ".").(string)
}

func push() {
	clientConfig, _ = cycle.GetClientConfig(keyExchanges, ciphers, cmdLine, "", hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, linuxMode)
	cycle.PushAndShow(clientConfig, &srcPath, &destPath)
}
