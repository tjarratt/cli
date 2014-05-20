package serviceauthtoken

import (
	"errors"
	"github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/models"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/tjarratt/cli/cf/terminal"
	cli "github.com/tjarratt/cg_cli"
)

type CreateServiceAuthTokenFields struct {
	ui            terminal.UI
	config        configuration.Reader
	authTokenRepo api.ServiceAuthTokenRepository
}

func NewCreateServiceAuthToken(ui terminal.UI, config configuration.Reader, authTokenRepo api.ServiceAuthTokenRepository) (cmd CreateServiceAuthTokenFields) {
	cmd.ui = ui
	cmd.config = config
	cmd.authTokenRepo = authTokenRepo
	return
}

func (cmd CreateServiceAuthTokenFields) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "create-service-auth-token",
		Description: "Create a service auth token",
		Usage:       "CF_NAME create-service-auth-token LABEL PROVIDER TOKEN",
	}
}

func (cmd CreateServiceAuthTokenFields) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 3 {
		err = errors.New("Incorrect usage")
		cmd.ui.FailWithUsage(c)
		return
	}

	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}
	return
}

func (cmd CreateServiceAuthTokenFields) Run(c *cli.Context) {
	cmd.ui.Say("Creating service auth token as %s...", terminal.EntityNameColor(cmd.config.Username()))

	serviceAuthTokenRepo := models.ServiceAuthTokenFields{
		Label:    c.Args()[0],
		Provider: c.Args()[1],
		Token:    c.Args()[2],
	}

	apiErr := cmd.authTokenRepo.Create(serviceAuthTokenRepo)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
}
