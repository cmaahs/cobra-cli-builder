package {{.Cmd}}

import (
	"{{.Program}}/config"

	"github.com/spf13/cobra"
)

var {{.Cmd}}Cmd = &cobra.Command{
	Use:   "{{.Cmd}}",
	Args:  cobra.MinimumNArgs(1),
	Short: "",
	Long: `EXAMPLES
	{{.Program}} {{.Cmd}} {{.Subcmd}}`,
	Run: func(cmd *cobra.Command, args []string) {},
}

var c *config.Config

func InitSubCommands(conf *config.Config) *cobra.Command {
	c = conf
	return {{.Cmd}}Cmd
}
