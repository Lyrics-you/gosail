package cli

import (
	"github.com/desertbit/grumble"
)

func init() {
	setCommand := &grumble.Command{
		Name: "set",
		Help: "Set the gosail config",
		Args: func(a *grumble.Args) {

		},
		Flags: func(f *grumble.Flags) {
			f.String("K", "key", "", "id_rsa.pub key filepath")
			f.String("", "ciphers", "", "ssh ciphers")
			f.String("", "keyexchanges", "", "ssh keyexchangesx")
			f.Int("N", "numlimit", 20, "max timeout")
			f.Int("T", "timelimit", 30, "max execute number")
		},
		Run: func(c *grumble.Context) error {
			setSetArgs(c)
			return nil
		},
	}
	Gosail.AddCommand(setCommand)
}

func setSetArgs(c *grumble.Context) {
	key = c.Flags.String("key")
	ciphers = c.Flags.String("ciphers")
	keyExchanges = c.Flags.String("keyexchanges")
	numLimit = c.Flags.Int("numlimit")
	timeLimit = c.Flags.Int("timelimit")
}
