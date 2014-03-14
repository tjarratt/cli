package user

import (
	"cf/api"
	"cf/configuration"
	"cf/errors"
	"cf/requirements"
	"cf/terminal"
	"github.com/tjarratt/cli"
)

type DeleteUserFields struct {
	ui       terminal.UI
	config   configuration.Reader
	userRepo api.UserRepository
}

func NewDeleteUser(ui terminal.UI, config configuration.Reader, userRepo api.UserRepository) (cmd DeleteUserFields) {
	cmd.ui = ui
	cmd.config = config
	cmd.userRepo = userRepo
	return
}

func (cmd DeleteUserFields) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 1 {
		err = errors.New("Invalid usage")
		cmd.ui.FailWithUsage(c, "delete-user")
		return
	}

	reqs = append(reqs, reqFactory.NewLoginRequirement())

	return
}

func (cmd DeleteUserFields) Run(c *cli.Context) {
	username := c.Args()[0]
	force := c.Bool("f")

	if !force && !cmd.ui.Confirm("Really delete user %s?%s",
		terminal.EntityNameColor(username),
		terminal.PromptColor(">"),
	) {
		return
	}

	cmd.ui.Say("Deleting user %s as %s...",
		terminal.EntityNameColor(username),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	user, apiErr := cmd.userRepo.FindByUsername(username)
	switch apiErr.(type) {
	case nil:
	case errors.ModelNotFoundError:
		cmd.ui.Ok()
		cmd.ui.Warn("User %s does not exist.", username)
		return
	default:
		cmd.ui.Failed(apiErr.Error())
		return
	}

	apiErr = cmd.userRepo.Delete(user.Guid)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
}
