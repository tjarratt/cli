package servicebroker

import (
	"github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/tjarratt/cli/cf/terminal"
	"github.com/codegangsta/cli"
)

type UpdateServiceBroker struct {
	ui     terminal.UI
	config configuration.Reader
	repo   api.ServiceBrokerRepository
}

func NewUpdateServiceBroker(ui terminal.UI, config configuration.Reader, repo api.ServiceBrokerRepository) (cmd UpdateServiceBroker) {
	cmd.ui = ui
	cmd.config = config
	cmd.repo = repo
	return
}

func (cmd UpdateServiceBroker) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "update-service-broker",
		Description: "Update a service broker",
		Usage:       "CF_NAME update-service-broker SERVICE_BROKER USERNAME PASSWORD URL",
	}
}

func (cmd UpdateServiceBroker) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 4 {
		cmd.ui.FailWithUsage(c)
	}

	reqs = append(reqs, requirementsFactory.NewLoginRequirement())

	return
}

func (cmd UpdateServiceBroker) Run(c *cli.Context) {
	serviceBroker, apiErr := cmd.repo.FindByName(c.Args()[0])
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Say("Updating service broker %s as %s...",
		terminal.EntityNameColor(serviceBroker.Name),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	serviceBroker.Username = c.Args()[1]
	serviceBroker.Password = c.Args()[2]
	serviceBroker.Url = c.Args()[3]

	apiErr = cmd.repo.Update(serviceBroker)

	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
}
