package commands

import (
	"github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/errors"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/tjarratt/cli/cf/terminal"
	cli "github.com/tjarratt/cg_cli"
)

type Password struct {
	ui      terminal.UI
	pwdRepo api.PasswordRepository
	config  configuration.ReadWriter
}

func NewPassword(ui terminal.UI, pwdRepo api.PasswordRepository, config configuration.ReadWriter) (cmd Password) {
	cmd.ui = ui
	cmd.pwdRepo = pwdRepo
	cmd.config = config
	return
}

func (cmd Password) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "passwd",
		ShortName:   "pw",
		Description: "Change user password",
		Usage:       "CF_NAME passwd",
	}
}

func (cmd Password) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	reqs = []requirements.Requirement{requirementsFactory.NewLoginRequirement()}
	return
}

func (cmd Password) Run(c *cli.Context) {
	oldPassword := cmd.ui.AskForPassword("Current Password")
	newPassword := cmd.ui.AskForPassword("New Password")
	verifiedPassword := cmd.ui.AskForPassword("Verify Password")

	if verifiedPassword != newPassword {
		cmd.ui.Failed("Password verification does not match")
		return
	}

	cmd.ui.Say("Changing password...")
	apiErr := cmd.pwdRepo.UpdatePassword(oldPassword, newPassword)

	switch typedErr := apiErr.(type) {
	case nil:
	case errors.HttpError:
		if typedErr.StatusCode() == 401 {
			cmd.ui.Failed("Current password did not match")
		} else {
			cmd.ui.Failed(apiErr.Error())
		}
	default:
		cmd.ui.Failed(apiErr.Error())
	}

	cmd.ui.Ok()
	cmd.config.ClearSession()
	cmd.ui.Say("Please log in again")
}
