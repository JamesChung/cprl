package console

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func IsGovCloud(cmd *cobra.Command, profile string) (bool, error) {
	gov, err := cmd.Flags().GetBool("gov-cloud")
	if err != nil {
		return gov, err
	}
	if gov {
		return gov, nil
	}
	gov = viper.GetBool(
		fmt.Sprintf(
			"%s.services.console.gov-cloud",
			profile,
		),
	)
	return gov, nil
}
