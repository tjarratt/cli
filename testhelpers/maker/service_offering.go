package maker

import "github.com/tjarratt/cli/cf/models"

func NewServiceOffering(label string) models.ServiceOffering {
	return models.ServiceOffering{ServiceOfferingFields: models.ServiceOfferingFields{
		Label:       label,
		Guid:        serviceOfferingGuid(),
		Description: "some service description",
	}}
}

var serviceOfferingGuid = guidGenerator("services")
