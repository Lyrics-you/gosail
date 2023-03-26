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
	version       bool
	username      string
	password      string
	key           string
	port          int
	ciphers       string
	keyExchanges  string
	config        string
	timeLimit     int
	numLimit      int
	linuxMode     bool
	jsonMode      bool
	selection     bool
	spinnerConfig *model.SpinConfig   = &model.SpinConfig{}
	clientConfig  *model.ClientConfig = &model.ClientConfig{}
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
			//todo: check``
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
	//User specified -> Load configuration -> Defualt value
	gosailConfig := cycle.GetGosailConfiguration()
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
		spinnerConfig.IsSelect = selection
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
		}
		// k8s
		if gosailConfig.K8s.Namespace != "" {
			namespace = gosailConfig.K8s.Namespace
		}
		if gosailConfig.K8s.AppName != "" {
			app = gosailConfig.K8s.AppName
		}
		if gosailConfig.K8s.Container != "" {
			container = gosailConfig.K8s.Container
		}
		if gosailConfig.K8s.Label != "" {
			label = gosailConfig.K8s.Label
		}
		if gosailConfig.K8s.Shell != "" {
			shell = gosailConfig.K8s.Shell
		}
		// mode
		if gosailConfig.Mode.JsonMode {
			jsonMode = gosailConfig.Mode.JsonMode
		}
		if gosailConfig.Mode.LinuxMode {
			linuxMode = gosailConfig.Mode.LinuxMode
		}
		if gosailConfig.Mode.Selection {
			selection = gosailConfig.Mode.Selection
		}
	}
	rootCmd.Execute()
}

func configExec() {
	clientConfig, err := cycle.MakeClientConfigFromJson(config)
	if err != nil {
		log.Error(err)
		return
	}
	spinnerConfig.IsSelect = selection
	sshResults := cycle.Exec(clientConfig, spinnerConfig)
	cycle.ShowExecResult(clientConfig.SshHosts, sshResults, &jsonMode, &linuxMode)
	if selection {
		client.LoginHostByID(clientConfig.SshHosts, sshResults, "")
	}
}
