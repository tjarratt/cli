package user

import (
	"github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/models"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/tjarratt/cli/cf/terminal"
	cli "github.com/tjarratt/cg_cli"
)

type SetOrgRole struct {
	ui       terminal.UI
	config   configuration.Reader
	userRepo api.UserRepository
	userReq  requirements.UserRequirement
	orgReq   requirements.OrganizationRequirement
}

func NewSetOrgRole(ui terminal.UI, config configuration.Reader, userRepo api.UserRepository) (cmd *SetOrgRole) {
	cmd = new(SetOrgRole)
	cmd.ui = ui
	cmd.config = config
	cmd.userRepo = userRepo
	return
}

func (cmd *SetOrgRole) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "set-org-role",
		Description: "Assign an org role to a user",
		Usage: "CF_NAME set-org-role USERNAME ORG ROLE\n\n" +
			"ROLES:\n" +
			"   OrgManager - Invite and manage users, select and change plans, and set spending limits\n" +
			"   BillingManager - Create and manage the billing account and payment info\n" +
			"   OrgAuditor - Read-only access to org info and reports\n",
	}
}

func (cmd *SetOrgRole) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
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

func (cmd *SetOrgRole) Run(c *cli.Context) {
	user := cmd.userReq.GetUser()
	org := cmd.orgReq.GetOrganization()
	role := models.UserInputToOrgRole[c.Args()[2]]

	cmd.ui.Say("Assigning role %s to user %s in org %s as %s...",
		terminal.EntityNameColor(role),
		terminal.EntityNameColor(user.Username),
		terminal.EntityNameColor(org.Name),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	apiErr := cmd.userRepo.SetOrgRole(user.Guid, org.Guid, role)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
}
