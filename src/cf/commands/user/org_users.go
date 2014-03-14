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

var orgRoles = []string{models.ORG_MANAGER, models.BILLING_MANAGER, models.ORG_AUDITOR}

var orgRoleToDisplayName = map[string]string{
	models.ORG_USER:        "USERS",
	models.ORG_MANAGER:     "ORG MANAGER",
	models.BILLING_MANAGER: "BILLING MANAGER",
	models.ORG_AUDITOR:     "ORG AUDITOR",
}

type OrgUsers struct {
	ui       terminal.UI
	config   configuration.Reader
	orgReq   requirements.OrganizationRequirement
	userRepo api.UserRepository
}

func NewOrgUsers(ui terminal.UI, config configuration.Reader, userRepo api.UserRepository) (cmd *OrgUsers) {
	cmd = new(OrgUsers)
	cmd.ui = ui
	cmd.config = config
	cmd.userRepo = userRepo
	return
}

func (cmd *OrgUsers) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 1 {
		err = errors.New("Incorrect usage")
		cmd.ui.FailWithUsage(c, "org-users")
		return
	}

	orgName := c.Args()[0]
	cmd.orgReq = reqFactory.NewOrganizationRequirement(orgName)
	reqs = append(reqs, reqFactory.NewLoginRequirement(), cmd.orgReq)

	return
}

func (cmd *OrgUsers) Run(c *cli.Context) {
	org := cmd.orgReq.GetOrganization()
	all := c.Bool("a")

	cmd.ui.Say("Getting users in org %s as %s...",
		terminal.EntityNameColor(org.Name),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	roles := orgRoles
	if all {
		roles = []string{models.ORG_USER}
	}

	for _, role := range roles {
		displayName := orgRoleToDisplayName[role]

		users, apiErr := cmd.userRepo.ListUsersInOrgForRole(org.Guid, role)

		cmd.ui.Say("")
		cmd.ui.Say("%s", terminal.HeaderColor(displayName))

		for _, user := range users {
			cmd.ui.Say("  %s", user.Username)
		}

		if apiErr != nil {
			cmd.ui.Failed("Failed fetching org-users for role %s.\n%s", apiErr.Error(), displayName)
			return
		}
	}
}
