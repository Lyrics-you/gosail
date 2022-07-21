package cli

import (
	"github.com/desertbit/grumble"
	"github.com/fatih/color"
	"gopkg.in/AlecAivazis/survey.v1"
)

var (
	black   = color.New(color.FgBlack, color.Bold).SprintFunc()
	red     = color.New(color.FgRed, color.Bold).SprintFunc()
	green   = color.New(color.FgGreen, color.Bold).SprintFunc()
	yellow  = color.New(color.FgYellow, color.Bold).SprintFunc()
	blue    = color.New(color.FgBlue, color.Bold).SprintFunc()
	magenta = color.New(color.FgMagenta, color.Bold).SprintFunc()
	// cyan    = color.New(color.FgCyan, color.Bold).SprintFunc()
	white = color.New(color.FgWhite, color.Bold).SprintFunc()
)
var qs = []*survey.Question{
	{
		Name: "color",
		Prompt: &survey.Select{
			Message: "Choose a color:",
			Options: []string{black("black"), red("red"), green("green"), yellow("yellow"), blue("blue"), magenta("magenta"), white("white")},
			Default: white("white"),
		},
	},
}

func init() {
	showCommand := &grumble.Command{
		Name: "color",
		Help: "Color for prompt",
		Args: func(a *grumble.Args) {

		},
		Flags: func(f *grumble.Flags) {

		},
		Run: func(_ *grumble.Context) error {
			ask()
			return nil
		},
	}
	Gosail.AddCommand(showCommand)
}

func ask() {
	// the answers will be written to this struct
	answers := struct {
		Color string `survey:"color"` // or you can tag fields to match a specific name
	}{}

	// perform the questions
	err := survey.Ask(qs, &answers)

	if err != nil {
		log.Error(err)
		return
	}

	switch answers.Color {
	case black("black"):
		promptColor = color.New(color.FgBlack)
	case red("red"):
		promptColor = color.New(color.FgRed)
	case green("green"):
		promptColor = color.New(color.FgGreen)
	case yellow("yellow"):
		promptColor = color.New(color.FgYellow)
	case blue("blue"):
		promptColor = color.New(color.FgBlue)
	case magenta("magenta"):
		promptColor = color.New(color.FgMagenta)
	case white("white"):
		promptColor = color.New(color.FgWhite, color.Bold)
	}
	Gosail.Config().PromptColor = promptColor
	Gosail.SetPrompt(prompt)
}
