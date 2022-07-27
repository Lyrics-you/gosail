package cli

import (
	"gosail/cycle"

	"github.com/desertbit/grumble"
)

var (
	srcPath  string
	destPath string
	tar      bool
	scp      bool
)

func init() {
	pullCommand := &grumble.Command{
		Aliases: []string{"download"},
		Name:    "pull",
		Help:    "Pull can copy file from hosts or pods, and create folders to distinguish",
		Args: func(a *grumble.Args) {
			a.String("src", "source path", grumble.Default(""))
			a.String("dest", "destination path", grumble.Default("."))
		},
		Flags: func(f *grumble.Flags) {
			f.Bool("", "tar", false, "tar pull's file")
			f.Bool("", "scp", false, "pull file by scp")
			f.String("", "src", "", "source path")
			f.String("", "dest", ".", "destination path")
		},
		Run: func(c *grumble.Context) error {
			setPullArgs(c)
			if isK8s {
				if scp {
					k8sPull()
				} else {
					k8sDownload()
				}
			} else {
				if scp {
					pull()
				} else {
					download()
				}
			}
			return nil
		},
	}
	Gosail.AddCommand(pullCommand)
}

func setPullArgs(c *grumble.Context) {
	srcPath = GetValue(c, "src", "").(string)
	destPath = GetValue(c, "dest", ".").(string)
	tar = c.Flags.Bool("tar")
	scp = c.Flags.Bool("scp")
	cmdLine = ""
}

func pull() {
	clientConfig, _ = cycle.GetClientConfig(keyExchanges, ciphers, cmdLine, "", hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, linuxMode)
	cycle.PullAndShow(clientConfig, &srcPath, &destPath, tar)
}

func download() {
	clientConfig, _ = cycle.GetClientConfig(keyExchanges, ciphers, cmdLine, "", hostLine, hostFile, ipLine, ipFile, username, password, key, port, numLimit, timeLimit, linuxMode)
	cycle.DownloadAndShow(clientConfig, &srcPath, &destPath, tar)
}
