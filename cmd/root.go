package cmd

import (
	"errors"
	"fmt"
	"gosail/cli"
	"gosail/client"
	"gosail/cycle"
	"gosail/model"

	"github.com/desertbit/grumble"
	"github.com/spf13/cobra"
)

var (
	version      bool
	username     string
	password     string
	key          string
	port         int
	ciphers      string
	keyExchanges string
	config       string
	timeLimit    int
	numLimit     int
	linuxMode    bool
	jsonMode     bool
	selection    bool
	clientConfig *model.ClientConfig = &model.ClientConfig{}
)

var (
	hostLine string
	hostFile string
	ipLine   string
	ipFile   string
)

var rootCmd = &cobra.Command{
	Use:   "gosail",
	Short: "gosail is a batch and concurrent command execution system.",
	Long: `
gosail is a free and open source batch and concurrent command execution system,
designed to execute commands on multiple servers or k8s pods and get results with speed and efficiency.
You can also copy(pull or push) files by it.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			Error(cmd, args, errors.New("unrecognized command"))
			return
		}
		if len(args) == 1 {
			cycle.LoginHost = args[0]
		}
	},
	PersistentPostRun: func(_ *cobra.Command, _ []string) {
		if config != "" {
			configExec()
			return
		}
		if !version {
			grumble.Main(cli.Gosail)
			fmt.Println("👌Finshed!")
		} else {
			showVersion()
		}
	},
}

func init() {
	// help
	rootCmd.PersistentFlags().BoolP("help", "?", false, "help for this command")
	// client
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "gosail version")
	rootCmd.PersistentFlags().StringVarP(&key, "key", "K", "", "id_rsa.pub key filepath")
	rootCmd.PersistentFlags().StringVarP(&ciphers, "ciphers", "", "", "ssh ciphers")
	rootCmd.PersistentFlags().StringVarP(&keyExchanges, "keyexchanges", "", "", "ssh keyexchanges")
	rootCmd.Flags().StringVarP(&config, "config", "", "", "host execute config")
	// limit
	rootCmd.PersistentFlags().IntVarP(&timeLimit, "timelimit", "T", 30, "max timeout")
	rootCmd.PersistentFlags().IntVarP(&numLimit, "numlimit", "N", 20, "max execute number")
	// cli
	rootCmd.Flags().StringVarP(&cycle.LoginHost, "hostfile", "", "", "for gosail cli loginhost")
	rootCmd.Flags().StringVarP(&cycle.LoginUser, "username", "u", "", "for gosail cli username")
	rootCmd.Flags().StringVarP(&cycle.LoginPwd, "password", "p", "", "for gosail cli password")
}

func Execute() {
	rootCmd.Execute()
}

func configExec() {
	clientConfig, err := cycle.MakeClientConfigFromJson(config)
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
