package user

import (
	cli "github.com/tjarratt/cg_cli"
	"github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/errors"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/tjarratt/cli/cf/terminal"
)

type DeleteUser struct {
	ui       terminal.UI
	config   configuration.Reader
	userRepo api.UserRepository
}

func NewDeleteUser(ui terminal.UI, config configuration.Reader, userRepo api.UserRepository) (cmd DeleteUser) {
	cmd.ui = ui
	cmd.config = config
	cmd.userRepo = userRepo
	return
}

func (cmd DeleteUser) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "delete-user",
		Description: "Delete a user",
		Usage:       "CF_NAME delete-user USERNAME [-f]",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "f", Usage: "Force deletion without confirmation"},
		},
	}
}

func (cmd DeleteUser) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 1 {
		err = errors.New("Invalid usage")
		cmd.ui.FailWithUsage(c)
		return
	}

	reqs = append(reqs, requirementsFactory.NewLoginRequirement())

	return
}

func (cmd DeleteUser) Run(c *cli.Context) {
	username := c.Args()[0]
	force := c.Bool("f")

	if !force && !cmd.ui.ConfirmDelete("user", username) {
		return
	}

	cmd.ui.Say("Deleting user %s as %s...",
		terminal.EntityNameColor(username),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	user, apiErr := cmd.userRepo.FindByUsername(username)
	switch apiErr.(type) {
	case nil:
	case *errors.ModelNotFoundError:
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
