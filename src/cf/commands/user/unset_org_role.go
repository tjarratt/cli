package user

import (
	"cf/api"
	"cf/configuration"
	"cf/models"
	"cf/requirements"
	"cf/terminal"
	"errors"
	"github.com/tjarratt/cli"
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

func (cmd *UnsetOrgRole) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 3 {
		err = errors.New("Incorrect Usage")
		cmd.ui.FailWithUsage(c, "unset-org-role")
		return
	}

	cmd.userReq = reqFactory.NewUserRequirement(c.Args()[0])
	cmd.orgReq = reqFactory.NewOrganizationRequirement(c.Args()[1])

	reqs = []requirements.Requirement{
		reqFactory.NewLoginRequirement(),
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
