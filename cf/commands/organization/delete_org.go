package organization

import (
	"github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/errors"
	"github.com/tjarratt/cli/cf/models"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/tjarratt/cli/cf/terminal"
	cli "github.com/tjarratt/cg_cli"
)

type DeleteOrg struct {
	ui      terminal.UI
	config  configuration.ReadWriter
	orgRepo api.OrganizationRepository
	orgReq  requirements.OrganizationRequirement
}

func NewDeleteOrg(ui terminal.UI, config configuration.ReadWriter, orgRepo api.OrganizationRepository) (cmd *DeleteOrg) {
	cmd = new(DeleteOrg)
	cmd.ui = ui
	cmd.config = config
	cmd.orgRepo = orgRepo
	return
}

func (cmd *DeleteOrg) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "delete-org",
		Description: "Delete an org",
		Usage:       "CF_NAME delete-org ORG [-f]",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "f", Usage: "Force deletion without confirmation"},
		},
	}
}

func (cmd *DeleteOrg) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 1 {
		cmd.ui.FailWithUsage(c)
	}

	reqs = []requirements.Requirement{requirementsFactory.NewLoginRequirement()}
	return
}

func (cmd *DeleteOrg) Run(c *cli.Context) {
	orgName := c.Args()[0]

	if !c.Bool("f") {
		if !cmd.ui.ConfirmDeleteWithAssociations("org", orgName) {
			return
		}
	}

	cmd.ui.Say("Deleting org %s as %s...",
		terminal.EntityNameColor(orgName),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	org, apiErr := cmd.orgRepo.FindByName(orgName)

	switch apiErr.(type) {
	case nil:
	case *errors.ModelNotFoundError:
		cmd.ui.Ok()
		cmd.ui.Warn("Org %s does not exist.", orgName)
		return
	default:
		cmd.ui.Failed(apiErr.Error())
		return
	}

	apiErr = cmd.orgRepo.Delete(org.Guid)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	if org.Guid == cmd.config.OrganizationFields().Guid {
		cmd.config.SetOrganizationFields(models.OrganizationFields{})
		cmd.config.SetSpaceFields(models.SpaceFields{})
	}

	cmd.ui.Ok()
	return
}
