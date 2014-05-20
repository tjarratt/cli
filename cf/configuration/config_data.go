package configuration

import (
	"github.com/tjarratt/cli/cf/models"
)

type AuthPromptType string

const (
	AuthPromptTypeText     AuthPromptType = "TEXT"
	AuthPromptTypePassword AuthPromptType = "PASSWORD"
)

type AuthPrompt struct {
	Type        AuthPromptType
	DisplayName string
}

type Data struct {
	ConfigVersion         int
	Target                string
	ApiVersion            string
	AuthorizationEndpoint string
	LoggregatorEndPoint   string
	UaaEndpoint           string
	AccessToken           string
	RefreshToken          string
	OrganizationFields    models.OrganizationFields
	SpaceFields           models.SpaceFields
	SSLDisabled           bool
	AsyncTimeout          uint
	Trace                 string
	ColorEnabled          string
}

func NewData() (data *Data) {
	data = new(Data)
	return
}
