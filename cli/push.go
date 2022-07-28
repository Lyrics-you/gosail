package cli

import (
	"gosail/cycle"
	"gosail/utils"

	"github.com/desertbit/grumble"
)

func init() {
	pushCommand := &grumble.Command{
		Aliases: []string{"upload"},
		Name:    "push",
		Help:    "Pull can copy file to hosts concurrently, and create folders that do not exist",
		Args: func(a *grumble.Args) {
			a.String("src", "source path", grumble.Default(""))
			a.String("dest", "estination path", grumble.Default("."))
		},
		Flags: func(f *grumble.Flags) {
			f.String("", "src", "", "source path")
			f.String("", "dest", ".", "destination path")
			f.Bool("", "scp", false, "push file by scp")
		},
		Run: func(c *grumble.Context) error {
			setPushArgs(c)
			if scp {
				push()
			} else {
				upload()
			}
			return nil
		},
	}
	Gosail.AddCommand(pushCommand)
}

func setPushArgs(c *grumble.Context) {
	srcPath = getValue(c, "src", "").(string)
	destPath = getValue(c, "dest", ".").(string)
	scp = c.Flags.Bool("scp")
}

func push() {
	clientConfig, err = cycle.GetClientConfig(keyExchanges, ciphers, cmdLine, "", hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, linuxMode)
	if err != nil && err != utils.ErrCmdListEmpty {
		log.Error(err)
		return
	}
	cycle.PushAndShow(clientConfig, &srcPath, &destPath)
}

func upload() {
	clientConfig, err = cycle.GetClientConfig(keyExchanges, ciphers, cmdLine, "", hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, linuxMode)
	if err != nil && err != utils.ErrCmdListEmpty {
		log.Error(err)
		return
	}
	cycle.UploadAndShow(clientConfig, &srcPath, &destPath)
}
