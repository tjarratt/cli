package servicebroker

import (
	cli "github.com/tjarratt/cg_cli"
	"github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/tjarratt/cli/cf/terminal"
)

type CreateServiceBroker struct {
	ui                terminal.UI
	config            configuration.Reader
	serviceBrokerRepo api.ServiceBrokerRepository
}

func NewCreateServiceBroker(ui terminal.UI, config configuration.Reader, serviceBrokerRepo api.ServiceBrokerRepository) (cmd CreateServiceBroker) {
	cmd.ui = ui
	cmd.config = config
	cmd.serviceBrokerRepo = serviceBrokerRepo
	return
}

func (cmd CreateServiceBroker) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "create-service-broker",
		Description: "Create a service broker",
		Usage:       "CF_NAME create-service-broker SERVICE_BROKER USERNAME PASSWORD URL",
	}
}

func (cmd CreateServiceBroker) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {

	if len(c.Args()) != 4 {
		cmd.ui.FailWithUsage(c)
	}

	reqs = append(reqs, requirementsFactory.NewLoginRequirement())

	return
}

func (cmd CreateServiceBroker) Run(c *cli.Context) {
	name := c.Args()[0]
	username := c.Args()[1]
	password := c.Args()[2]
	url := c.Args()[3]

	cmd.ui.Say("Creating service broker %s as %s...",
		terminal.EntityNameColor(name),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	apiErr := cmd.serviceBrokerRepo.Create(name, url, username, password)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
}
