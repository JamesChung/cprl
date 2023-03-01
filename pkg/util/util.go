package util

import (
	"fmt"
	"os"
	"strings"

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
		fmt.Fprintf(os.Stderr, "%s\n", err)
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
