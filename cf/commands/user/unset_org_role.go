package user

import (
	cli "github.com/tjarratt/cg_cli"
	"github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/models"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/tjarratt/cli/cf/terminal"
)

type UnsetOrgRole struct {
	ui       terminal.UI
	config   configuration.Reader
	userRepo api.UserRepository
	userReq  requirements.UserRequirement
	orgReq   requirements.OrganizationRequirement
}

func NewUnsetOrgRole(ui terminal.UI, config configuration.Reader, userRepo api.UserRepository) (cmd *UnsetOrgRole) {
	cmd = new(UnsetOrgRole)
	cmd.ui = ui
	cmd.config = config
	cmd.userRepo = userRepo

	return
}

func (cmd *UnsetOrgRole) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "unset-org-role",
		Description: "Remove an org role from a user",
		Usage: "CF_NAME unset-org-role USERNAME ORG ROLE\n\n" +
			"ROLES:\n" +
			"   OrgManager - Invite and manage users, select and change plans, and set spending limits\n" +
			"   BillingManager - Create and manage the billing account and payment info\n" +
			"   OrgAuditor - Read-only access to org info and reports\n",
	}
}

func (cmd *UnsetOrgRole) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 3 {
		cmd.ui.FailWithUsage(c)
	}

	cmd.userReq = requirementsFactory.NewUserRequirement(c.Args()[0])
	cmd.orgReq = requirementsFactory.NewOrganizationRequirement(c.Args()[1])

	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		cmd.userReq,
		cmd.orgReq,
	}

	return
}

func (cmd *UnsetOrgRole) Run(c *cli.Context) {
	role := models.UserInputToOrgRole[c.Args()[2]]
	user := cmd.userReq.GetUser()
	org := cmd.orgReq.GetOrganization()

	cmd.ui.Say("Removing role %s from user %s in org %s as %s...",
		terminal.EntityNameColor(role),
		terminal.EntityNameColor(c.Args()[0]),
		terminal.EntityNameColor(c.Args()[1]),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	apiErr := cmd.userRepo.UnsetOrgRole(user.Guid, org.Guid, role)

	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
}
