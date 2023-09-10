package {{.Cmd}}

import (
	"encoding/json"

	"{{.Program}}/objects"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// {{.Subcmd}}Cmd represents the {{.Subcmd}} command
var {{.Subcmd}}Cmd = &cobra.Command{
	Use:   "{{.Subcmd}}",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := DoThis()
		if err != nil {
			logrus.WithError(err).Error("Failed to express the version")
		}
		if !c.FormatOverridden {
			c.OutputFormat = "text"
		}
		c.OutputData(&out)
	},
}

func DoThis() (objects.{{.SubcmdPublic}}, error) {
	var {{.Subcmd}}Data objects.{{.SubcmdPublic}}
	err := json.Unmarshal([]byte({{.SubcmdPublic}}JSON), &{{.Subcmd}}Data)
	if err != nil {
		return {{.Subcmd}}Data, errors.Wrap(err,"Failed to unmarshal JSON")
	}

	return {{.Subcmd}}Data, nil
}

func init() {
	{{.Cmd}}Cmd.AddCommand({{.Subcmd}}Cmd)
}
