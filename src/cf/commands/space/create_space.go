package space

import (
	"cf"
	"cf/api"
	"cf/commands/user"
	"cf/configuration"
	"cf/errors"
	"cf/models"
	"cf/requirements"
	"cf/terminal"
	"github.com/tjarratt/cli"
)

type CreateSpace struct {
	ui              terminal.UI
	config          configuration.Reader
	spaceRepo       api.SpaceRepository
	orgRepo         api.OrganizationRepository
	userRepo        api.UserRepository
	spaceRoleSetter user.SpaceRoleSetter
}

func NewCreateSpace(ui terminal.UI, config configuration.Reader, spaceRoleSetter user.SpaceRoleSetter, spaceRepo api.SpaceRepository, orgRepo api.OrganizationRepository, userRepo api.UserRepository) (cmd CreateSpace) {
	cmd.ui = ui
	cmd.config = config
	cmd.spaceRoleSetter = spaceRoleSetter
	cmd.spaceRepo = spaceRepo
	cmd.orgRepo = orgRepo
	cmd.userRepo = userRepo
	return
}

func (cmd CreateSpace) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) == 0 {
		err = errors.New("Incorrect Usage")
		cmd.ui.FailWithUsage(c, "create-space")
		return
	}

	reqs = []requirements.Requirement{reqFactory.NewLoginRequirement()}
	if c.String("o") == "" {
		reqs = append(reqs, reqFactory.NewTargetedOrgRequirement())
	}

	return
}

func (cmd CreateSpace) Run(c *cli.Context) {
	spaceName := c.Args()[0]
	orgName := c.String("o")
	orgGuid := ""
	if orgName == "" {
		orgName = cmd.config.OrganizationFields().Name
		orgGuid = cmd.config.OrganizationFields().Guid
	}

	cmd.ui.Say("Creating space %s in org %s as %s...",
		terminal.EntityNameColor(spaceName),
		terminal.EntityNameColor(orgName),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	if orgGuid == "" {
		org, apiErr := cmd.orgRepo.FindByName(orgName)
		switch apiErr.(type) {
		case nil:
		case errors.ModelNotFoundError:
			cmd.ui.Failed("Org %s does not exist or is not accessible", orgName)
			return
		default:
			cmd.ui.Failed("Error finding org %s\n%s", orgName, apiErr.Error())
			return
		}

		orgGuid = org.Guid
	}

	space, apiErr := cmd.spaceRepo.Create(spaceName, orgGuid)
	if apiErr != nil {
		if apiErr.ErrorCode() == cf.SPACE_EXISTS {
			cmd.ui.Ok()
			cmd.ui.Warn("Space %s already exists", spaceName)
			return
		}
		cmd.ui.Failed(apiErr.Error())
		return
	}
	cmd.ui.Ok()

	var err error

	err = cmd.spaceRoleSetter.SetSpaceRole(space, models.SPACE_MANAGER, cmd.config.UserGuid(), cmd.config.Username())
	if err != nil {
		cmd.ui.Failed(err.Error())
		return
	}

	err = cmd.spaceRoleSetter.SetSpaceRole(space, models.SPACE_DEVELOPER, cmd.config.UserGuid(), cmd.config.Username())
	if err != nil {
		cmd.ui.Failed(err.Error())
		return
	}

	cmd.ui.Say("\nTIP: Use '%s' to target new space", terminal.CommandColor(cf.Name()+" target -o "+orgName+" -s "+space.Name))
}
