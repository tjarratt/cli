package commands

import (
	cli "github.com/tjarratt/cg_cli"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/help"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/tjarratt/cli/cf/terminal"
)

type Help struct {
	ui terminal.UI
}

func NewHelp(ui terminal.UI) Help {
	return Help{ui: ui}
}

func (cmd Help) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "help",
		Description: "helps with a bunch of stuff, right?",
		Usage:       "CF_NAME help [COMMAND_NAME]",
		Flags:       []cli.Flag{},
	}
}

func (cmd Help) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	return
}

func (cmd Help) Run(c *cli.Context) {
	args := c.Args()
	if len(args) > 0 {
		cli.ShowCommandHelp(c, args[0])
	} else {
		help.ShowAppHelp(help.AppHelpTemplate(), c.App)
	}
}
