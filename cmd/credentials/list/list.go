package list

import (
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/JamesChung/cprl/pkg/util"
)

var (
	shortMessage = "lists AWS profiles"

	example = templates.Examples(`cprl credentials list`)
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l", "li"},
		Short:   shortMessage,
		Example: example,
		Run:     runCmd,
	}
	return cmd
}

func generateLeveledListItems(items []string, level int) []pterm.LeveledListItem {
	leveledListItems := make([]pterm.LeveledListItem, 0, len(items))
	for _, item := range items {
		leveledListItems = append(leveledListItems, pterm.LeveledListItem{
			Level: level,
			Text:  item,
		})
	}
	return leveledListItems
}

func runCmd(cmd *cobra.Command, args []string) {
	configProfiles, err := util.GetConfigProfiles()
	util.ExitOnErr(err)
	credentialsProfiles, err := util.GetCredentialsProfiles()
	util.ExitOnErr(err)

	// Config
	ll := pterm.LeveledList{pterm.LeveledListItem{Level: 0, Text: "Config"}}
	ll = append(ll, generateLeveledListItems(configProfiles, 1)...)

	// Credentials
	ll = append(ll, pterm.LeveledListItem{Level: 0, Text: "Credentials"})
	ll = append(ll, generateLeveledListItems(credentialsProfiles, 1)...)

	// Render tree
	pterm.DefaultTree.WithRoot(putils.TreeFromLeveledList(ll)).Render()
}
