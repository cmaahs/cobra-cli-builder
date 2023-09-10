package new

import (
	"new-cli-subcmd/config"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Args:  cobra.MinimumNArgs(1),
	Short: "",
	Long: `EXAMPLES
	new-cli-subcmd new subcmd`,
	Run: func(cmd *cobra.Command, args []string) {},
}

var c *config.Config

func InitSubCommands(conf *config.Config) *cobra.Command {
	c = conf
	return newCmd
}
