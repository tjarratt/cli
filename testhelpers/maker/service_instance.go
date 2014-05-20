package maker

import "github.com/tjarratt/cli/cf/models"

var serviceInstanceGuid func() string

func init() {
	serviceInstanceGuid = guidGenerator("services")
}

func NewServiceInstance(name string) (service models.ServiceInstance) {
	service.Name = name
	service.Guid = serviceInstanceGuid()
	return
}
