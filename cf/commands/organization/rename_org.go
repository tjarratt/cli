package organization

import (
	cli "github.com/tjarratt/cg_cli"
	"github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/tjarratt/cli/cf/terminal"
)

type RenameOrg struct {
	ui      terminal.UI
	config  configuration.ReadWriter
	orgRepo api.OrganizationRepository
	orgReq  requirements.OrganizationRequirement
}

func NewRenameOrg(ui terminal.UI, config configuration.ReadWriter, orgRepo api.OrganizationRepository) (cmd *RenameOrg) {
	cmd = new(RenameOrg)
	cmd.ui = ui
	cmd.config = config
	cmd.orgRepo = orgRepo
	return
}

func (cmd *RenameOrg) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "rename-org",
		Description: "Rename an org",
		Usage:       "CF_NAME rename-org ORG NEW_ORG",
	}
}

func (cmd *RenameOrg) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 2 {
		cmd.ui.FailWithUsage(c)
	}

	cmd.orgReq = requirementsFactory.NewOrganizationRequirement(c.Args()[0])
	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		cmd.orgReq,
	}
	return
}

func (cmd *RenameOrg) Run(c *cli.Context) {
	org := cmd.orgReq.GetOrganization()
	newName := c.Args()[1]

	cmd.ui.Say("Renaming org %s to %s as %s...",
		terminal.EntityNameColor(org.Name),
		terminal.EntityNameColor(newName),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	apiErr := cmd.orgRepo.Rename(org.Guid, newName)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}
	cmd.ui.Ok()

	if org.Guid == cmd.config.OrganizationFields().Guid {
		org.Name = newName
		cmd.config.SetOrganizationFields(org.OrganizationFields)
	}
}
