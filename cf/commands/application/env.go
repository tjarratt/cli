package application

import (
	"strings"

	cli "github.com/tjarratt/cg_cli"
	"github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/errors"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/tjarratt/cli/cf/terminal"
)

type Env struct {
	ui      terminal.UI
	config  configuration.Reader
	appRepo api.ApplicationRepository
}

func NewEnv(ui terminal.UI, config configuration.Reader, appRepo api.ApplicationRepository) (cmd *Env) {
	cmd = new(Env)
	cmd.ui = ui
	cmd.config = config
	cmd.appRepo = appRepo
	return
}

func (cmd *Env) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "env",
		ShortName:   "e",
		Description: "Show all env variables for an app",
		Usage:       "CF_NAME env APP",
	}
}

func (cmd *Env) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) ([]requirements.Requirement, error) {
	if len(c.Args()) != 1 {
		cmd.ui.FailWithUsage(c)
	}

	return []requirements.Requirement{requirementsFactory.NewLoginRequirement()}, nil
}

func (cmd *Env) Run(c *cli.Context) {
	app, err := cmd.appRepo.Read(c.Args()[0])
	if notFound, ok := err.(*errors.ModelNotFoundError); ok {
		cmd.ui.Failed(notFound.Error())
	}

	cmd.ui.Say("Getting env variables for app %s in org %s / space %s as %s...",
		terminal.EntityNameColor(app.Name),
		terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
		terminal.EntityNameColor(cmd.config.SpaceFields().Name),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	envVars, vcapServices, err := cmd.appRepo.ReadEnv(app.Guid)
	if err != nil {
		cmd.ui.Failed(err.Error())
	}

	cmd.ui.Ok()
	cmd.ui.Say("")

	if len(vcapServices) > 0 {
		cmd.ui.Say("System-Provided:")
		for _, line := range strings.Split(vcapServices, "\n") {
			cmd.ui.Say(line)
		}
	} else {
		cmd.ui.Say("No system-provided env variables have been set")
	}

	if len(envVars) == 0 {
		cmd.ui.Say("No user-defined env variables have been set")
		return
	}

	cmd.ui.Say("User-Provided:")
	for key, value := range envVars {
		cmd.ui.Say("%s: %s", key, terminal.EntityNameColor(value))
	}
}
