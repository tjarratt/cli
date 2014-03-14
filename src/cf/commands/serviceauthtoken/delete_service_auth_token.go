package serviceauthtoken

import (
	"cf/api"
	"cf/configuration"
	"cf/errors"
	"cf/requirements"
	"cf/terminal"
	"fmt"
	"github.com/tjarratt/cli"
)

type DeleteServiceAuthTokenFields struct {
	ui            terminal.UI
	config        configuration.Reader
	authTokenRepo api.ServiceAuthTokenRepository
}

func NewDeleteServiceAuthToken(ui terminal.UI, config configuration.Reader, authTokenRepo api.ServiceAuthTokenRepository) (cmd DeleteServiceAuthTokenFields) {
	cmd.ui = ui
	cmd.config = config
	cmd.authTokenRepo = authTokenRepo
	return
}

func (cmd DeleteServiceAuthTokenFields) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 2 {
		err = errors.New("Incorrect usage")
		cmd.ui.FailWithUsage(c, "delete-service-auth-token")
		return
	}

	reqs = append(reqs, reqFactory.NewLoginRequirement())
	return
}

func (cmd DeleteServiceAuthTokenFields) Run(c *cli.Context) {
	tokenLabel := c.Args()[0]
	tokenProvider := c.Args()[1]

	if c.Bool("f") == false {
		response := cmd.ui.Confirm(
			"Are you sure you want to delete %s?%s",
			terminal.EntityNameColor(fmt.Sprintf("%s %s", tokenLabel, tokenProvider)),
			terminal.PromptColor(">"),
		)
		if response == false {
			return
		}
	}

	cmd.ui.Say("Deleting service auth token as %s", terminal.EntityNameColor(cmd.config.Username()))
	token, apiErr := cmd.authTokenRepo.FindByLabelAndProvider(tokenLabel, tokenProvider)

	switch apiErr.(type) {
	case nil:
	case errors.ModelNotFoundError:
		cmd.ui.Ok()
		cmd.ui.Warn("Service Auth Token %s %s does not exist.", tokenLabel, tokenProvider)
		return
	default:
		cmd.ui.Failed(apiErr.Error())
	}

	apiErr = cmd.authTokenRepo.Delete(token)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
}
