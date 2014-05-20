package configuration

import (
	"encoding/json"
	"github.com/tjarratt/cli/cf/models"
)

type configJsonV3 struct {
	ConfigVersion         int
	Target                string
	ApiVersion            string
	AuthorizationEndpoint string
	LoggregatorEndpoint   string
	UaaEndpoint           string
	AccessToken           string
	RefreshToken          string
	OrganizationFields    models.OrganizationFields
	SpaceFields           models.SpaceFields
	SSLDisabled           bool
	AsyncTimeout          uint
	Trace                 string
	ColorEnabled          string // need to be able to express true, false and undefined
}

func JsonMarshalV3(config *Data) (output []byte, err error) {
	return json.Marshal(configJsonV3{
		ConfigVersion:         3,
		Target:                config.Target,
		ApiVersion:            config.ApiVersion,
		AuthorizationEndpoint: config.AuthorizationEndpoint,
		LoggregatorEndpoint:   config.LoggregatorEndPoint,
		UaaEndpoint:           config.UaaEndpoint,
		AccessToken:           config.AccessToken,
		RefreshToken:          config.RefreshToken,
		OrganizationFields:    config.OrganizationFields,
		SpaceFields:           config.SpaceFields,
		SSLDisabled:           config.SSLDisabled,
		Trace:                 config.Trace,
		AsyncTimeout:          config.AsyncTimeout,
		ColorEnabled:          config.ColorEnabled,
	})
}

func JsonUnmarshalV3(input []byte, config *Data) (err error) {
	configJson := new(configJsonV3)

	err = json.Unmarshal(input, configJson)
	if err != nil {
		return
	}

	if configJson.ConfigVersion != 3 {
		return
	}

	config.Target = configJson.Target
	config.ApiVersion = configJson.ApiVersion
	config.AccessToken = configJson.AccessToken
	config.RefreshToken = configJson.RefreshToken
	config.SpaceFields = configJson.SpaceFields
	config.OrganizationFields = configJson.OrganizationFields
	config.LoggregatorEndPoint = configJson.LoggregatorEndpoint
	config.AuthorizationEndpoint = configJson.AuthorizationEndpoint
	config.UaaEndpoint = configJson.UaaEndpoint
	config.SSLDisabled = configJson.SSLDisabled
	config.AsyncTimeout = configJson.AsyncTimeout
	config.Trace = configJson.Trace
	config.ColorEnabled = configJson.ColorEnabled

	return
}
