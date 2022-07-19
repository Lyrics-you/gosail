package cmd

import (
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login host to do something",
	Long: `
eg. : gosail login -h <hostfile> [-u "<username>"] [-p "<password>"] [--prot "<port>"]

If the ssh port is 22, can omit port arg
eg. : gosail login -h <hostfile> [-u "<username>"] [-p "<password>"]

If the hostfile or hostline contain hosts in the format username@host, can omit u arg
eg. : gosail login -h <hostfile> [-p "<password>"]

If specified the K arg or has default id_rsa.pub key, can omit p arg
eg. : gosail login -h <hostfile>
`,
	// TraverseChildren: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
	},
	Run: func(_ *cobra.Command, args []string) {
		if len(args) > 0 {
			// cmd.Error(cmd, args, errors.New("unrecognized command"))
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	//login
	loginCmd.PersistentFlags().StringVarP(&hostLine, "host", "", "", "hostline")
	loginCmd.PersistentFlags().StringVarP(&hostFile, "hostfile", "h", "", "hostfile")
	loginCmd.PersistentFlags().StringVarP(&ipLine, "ip", "", "", "ipline")
	loginCmd.PersistentFlags().StringVarP(&ipFile, "ipfile", "i", "", "ipfile")

	loginCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "host username")
	loginCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "host password")
	loginCmd.PersistentFlags().IntVarP(&port, "port", "", 22, "ssh port")
}
