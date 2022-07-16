package cli

import (
	"github.com/desertbit/grumble"
)

var (
	jsonMode  bool
	linuxMode = false
	selection bool
)

func init() {
	modeCommand := &grumble.Command{
		Name: "mode",
		Help: "Mode offers choices of exec output formats",
		Args: func(a *grumble.Args) {

		},
		Flags: func(f *grumble.Flags) {
			f.Bool("l", "linuxmode", false, "linux mode")
			f.Bool("j", "jsonmode", false, "json mode")
			f.Bool("s", "selection", false, "selection")
		},
		Run: func(c *grumble.Context) error {
			setModeArgs(c)
			return nil
		},
	}
	Gosail.AddCommand(modeCommand)
}

func setModeArgs(c *grumble.Context) {
	linuxMode = c.Flags.Bool("linuxmode")
	jsonMode = c.Flags.Bool("jsonmode")
	selection = c.Flags.Bool("selection")
	if isK8s {
		linuxMode = true
	}
}
