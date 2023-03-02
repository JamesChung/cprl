package util

import (
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func AddGroup(parent *cobra.Command, title string, cmds ...*cobra.Command) {
	group := &cobra.Group{
		Title: title,
		ID:    title,
	}
	parent.AddGroup(group)
	for _, cmd := range cmds {
		cmd.GroupID = group.ID
		parent.AddCommand(cmd)
	}
}

func Basename(str string) string {
	s := strings.Split(str, "/")
	return s[len(s)-1]
}

func ExitOnErr(err error) {
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
}

func GetFlagString(cmd *cobra.Command, str string) (string, error) {
	val, err := cmd.Flags().GetString(str)
	if err != nil {
		return "", err
	}
	return val, nil
}

func GetFlagBool(cmd *cobra.Command, str string) (bool, error) {
	val, err := cmd.Flags().GetBool(str)
	if err != nil {
		return false, err
	}
	return val, nil
}

func Spinner(startMsg string, closure func()) {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Color("blue", "bold")
	s.Suffix = pterm.Blue(" ", startMsg)
	s.Start()
	closure()
	s.Stop()
}

func SpinnerWithStatusMsg(startMsg string, closure func() (string, error)) {
	var str string
	var err error
	Spinner(startMsg, func() { str, err = closure() })
	switch {
	case err != nil:
		pterm.Error.Println(err)
		return
	case str != "":
		pterm.Success.Println(str)
		return
	}
	pterm.Success.Println()
}
