package cprlerrs

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var (
	MissingDefaultProfileErr = errors.New("cprl.yaml is missing [default] profile")
	MissingRepositoriesErr   = errors.New("cprl.yaml is missing [repositories]")
)

func ExitOnErr(cmd *cobra.Command, err error) {
	if err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}
}
