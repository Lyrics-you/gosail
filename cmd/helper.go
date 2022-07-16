package cmd

import (
	"gosail/logger"
	"strings"

	"github.com/spf13/cobra"
)

var (
	log = logger.Logger()
)

func Error(cmd *cobra.Command, args []string, err error) {
	log.Errorf("execute 'gosail %s %s' error, %v", cmd.Name(), strings.Join(args, " "), err)
}
