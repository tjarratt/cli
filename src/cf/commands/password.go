package commands

import (
	"cf/api"
	"cf/configuration"
	"cf/errors"
	"cf/requirements"
	"cf/terminal"
	"github.com/tjarratt/cli"
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

func (cmd Password) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	reqs = []requirements.Requirement{reqFactory.NewLoginRequirement()}
	return
}

func (cmd Password) Run(c *cli.Context) {
	oldPassword := cmd.ui.AskForPassword("Current Password%s", terminal.PromptColor(">"))
	newPassword := cmd.ui.AskForPassword("New Password%s", terminal.PromptColor(">"))
	verifiedPassword := cmd.ui.AskForPassword("Verify Password%s", terminal.PromptColor(">"))

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
