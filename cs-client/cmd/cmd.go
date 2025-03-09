package cmd

import (
	"fmt"
	"strings"

	"github.com/ignite/apps/cs-client/gen"
	"github.com/ignite/cli/v29/ignite/services/plugin"
)

// GetCommands returns the list of spinner app commands.
func GetCommands() []*plugin.Command {
	return []*plugin.Command{
		{
			Use:   "cs-client",
			Short: "Generates csharp client",
			Long:  "Generates csharp client",
			Flags: []*plugin.Flag{
				{
					Name:      "yes",
					Shorthand: "y",
					Type:      plugin.FlagTypeBool,
					Usage:     "answers interactive yes/no questions with yes",
				},
				{
					Name:      "out",
					Shorthand: "o",
					Type:      plugin.FlagTypeString,
					Usage:     "csharp output directory",
				},
				{
					Name: "components",
					Type: plugin.FlagTypeStringSlice,
					Usage: fmt.Sprintf(
						"components to be generated; options: [%s]",
						strings.Join(gen.Component_stringValues(), ", "),
					),
				},
			},
			PlaceCommandUnder: "generate",
		},
	}
}
