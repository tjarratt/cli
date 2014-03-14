package service

import (
	"cf/api"
	"cf/configuration"
	"cf/errors"
	"cf/requirements"
	"cf/terminal"
	"github.com/tjarratt/cli"
)

type PurgeServiceOffering struct {
	ui          terminal.UI
	serviceRepo api.ServiceRepository
}

func (cmd PurgeServiceOffering) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 1 {
		err = errors.New("incorrect usage")
		cmd.ui.FailWithUsage(c, "purge-service-offering")
	}
	return
}

func (cmd PurgeServiceOffering) Run(c *cli.Context) {
	serviceName := c.Args()[0]

	offering, apiErr := cmd.serviceRepo.FindServiceOfferingByLabelAndProvider(serviceName, c.String("p"))

	switch apiErr.(type) {
	case nil:
	case errors.ModelNotFoundError:
		cmd.ui.Warn("Service offering does not exist\nTIP: If you are trying to purge a v1 service offering, you must set the -p flag.")
		return
	default:
		cmd.ui.Failed(apiErr.Error())
	}

	confirmed := c.Bool("f")
	if !confirmed {
		cmd.ui.Warn(`Warning: This operation assumes that the service broker responsible for this service offering is no longer available, and all service instances have been deleted, leaving orphan records in Cloud Foundry's database. All knowledge of the service will be removed from Cloud Foundry, including service instances and service bindings. No attempt will be made to contact the service broker; running this command without destroying the service broker will cause orphan service instances. After running this command you may want to run either delete-service-auth-token or delete-service-broker to complete the cleanup.`)
		confirmed = cmd.ui.Confirm("Really purge service offering %s from Cloud Foundry?", serviceName)
	}

	if !confirmed {
		return
	}
	cmd.ui.Say("Purging service %s...", serviceName)
	cmd.serviceRepo.PurgeServiceOffering(offering)
	cmd.ui.Ok()

	return
}

func NewPurgeServiceOffering(ui terminal.UI, config configuration.Reader, serviceRepo api.ServiceRepository) (cmd PurgeServiceOffering) {
	cmd.ui = ui
	cmd.serviceRepo = serviceRepo
	return
}
