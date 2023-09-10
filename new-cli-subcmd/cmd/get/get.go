package get

import (
	"new-cli-subcmd/config"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Args:  cobra.MinimumNArgs(1),
	Short: "",
	Long: `EXAMPLES
	new-cli-subcmd get version`,
	Run: func(cmd *cobra.Command, args []string) {},
}

var c *config.Config

func InitSubCommands(conf *config.Config) *cobra.Command {
	c = conf
	return getCmd
}
