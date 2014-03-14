package commands

import (
	"cf"
	"cf/api"
	"cf/configuration"
	"cf/requirements"
	"cf/terminal"
	"errors"
	"github.com/tjarratt/cli"
)

type Authenticate struct {
	ui            terminal.UI
	config        configuration.ReadWriter
	authenticator api.AuthenticationRepository
}

func NewAuthenticate(ui terminal.UI, config configuration.ReadWriter, authenticator api.AuthenticationRepository) (cmd Authenticate) {
	cmd.ui = ui
	cmd.config = config
	cmd.authenticator = authenticator
	return
}

func (cmd Authenticate) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 2 {
		err = errors.New("Incorrect Usage")
		cmd.ui.FailWithUsage(c, "auth")
		return
	}

	reqs = append(reqs, reqFactory.NewApiEndpointRequirement())
	return
}

func (cmd Authenticate) Run(c *cli.Context) {
	cmd.config.ClearSession()

	cmd.ui.Say("API endpoint: %s", terminal.EntityNameColor(cmd.config.ApiEndpoint()))
	cmd.ui.Say("Authenticating...")

	apiErr := cmd.authenticator.Authenticate(map[string]string{
		"username": c.Args()[0],
		"password": c.Args()[1],
	})
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
	cmd.ui.Say("Use '%s' to view or set your target org and space", terminal.CommandColor(cf.Name()+" target"))
	return
}
