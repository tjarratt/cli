package requirements

import (
	"github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/models"
	"github.com/tjarratt/cli/cf/terminal"
)

type ServiceInstanceRequirement interface {
	Requirement
	GetServiceInstance() models.ServiceInstance
}

type serviceInstanceApiRequirement struct {
	name            string
	ui              terminal.UI
	serviceRepo     api.ServiceRepository
	serviceInstance models.ServiceInstance
}

func NewServiceInstanceRequirement(name string, ui terminal.UI, sR api.ServiceRepository) (req *serviceInstanceApiRequirement) {
	req = new(serviceInstanceApiRequirement)
	req.name = name
	req.ui = ui
	req.serviceRepo = sR
	return
}

func (req *serviceInstanceApiRequirement) Execute() (success bool) {
	var apiErr error
	req.serviceInstance, apiErr = req.serviceRepo.FindInstanceByName(req.name)

	if apiErr != nil {
		req.ui.Failed(apiErr.Error())
		return false
	}

	return true
}

func (req *serviceInstanceApiRequirement) GetServiceInstance() models.ServiceInstance {
	return req.serviceInstance
}
