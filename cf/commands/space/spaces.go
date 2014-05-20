package space

import (
	"github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/models"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/tjarratt/cli/cf/terminal"
	cli "github.com/tjarratt/cg_cli"
)

type ListSpaces struct {
	ui        terminal.UI
	config    configuration.Reader
	spaceRepo api.SpaceRepository
}

func NewListSpaces(ui terminal.UI, config configuration.Reader, spaceRepo api.SpaceRepository) (cmd ListSpaces) {
	cmd.ui = ui
	cmd.config = config
	cmd.spaceRepo = spaceRepo
	return
}

func (cmd ListSpaces) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "spaces",
		Description: "List all spaces in an org",
		Usage:       "CF_NAME spaces",
	}
}

func (cmd ListSpaces) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		requirementsFactory.NewTargetedOrgRequirement(),
	}
	return
}

func (cmd ListSpaces) Run(c *cli.Context) {
	cmd.ui.Say("Getting spaces in org %s as %s...\n",
		terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
		terminal.EntityNameColor(cmd.config.Username()))

	foundSpaces := false
	table := cmd.ui.Table([]string{"name"})
	apiErr := cmd.spaceRepo.ListSpaces(func(space models.Space) bool {
		table.Add([]string{space.Name})
		foundSpaces = true
		return true
	})
	table.Print()

	if apiErr != nil {
		cmd.ui.Failed("Failed fetching spaces.\n%s", apiErr.Error())
		return
	}

	if !foundSpaces {
		cmd.ui.Say("No spaces found")
	}
}
