package commands

import (
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/tjarratt/cli/cf/terminal"
	"github.com/codegangsta/cli"
)

type Logout struct {
	ui     terminal.UI
	config configuration.ReadWriter
}

func NewLogout(ui terminal.UI, config configuration.ReadWriter) (cmd Logout) {
	cmd.ui = ui
	cmd.config = config
	return
}

func (cmd Logout) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "logout",
		ShortName:   "lo",
		Description: "Log user out",
		Usage:       "CF_NAME logout",
	}
}

func (cmd Logout) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	return
}

func (cmd Logout) Run(c *cli.Context) {
	cmd.ui.Say("Logging out...")
	cmd.config.ClearSession()
	cmd.ui.Ok()
}
