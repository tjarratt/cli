package api

import (
	"fmt"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/errors"
	"github.com/tjarratt/cli/cf/net"
	"strings"
)

type PasswordRepository interface {
	UpdatePassword(old string, new string) error
}

type CloudControllerPasswordRepository struct {
	config  configuration.Reader
	gateway net.Gateway
}

func NewCloudControllerPasswordRepository(config configuration.Reader, gateway net.Gateway) (repo CloudControllerPasswordRepository) {
	repo.config = config
	repo.gateway = gateway
	return
}

func (repo CloudControllerPasswordRepository) UpdatePassword(old string, new string) error {
	uaaEndpoint := repo.config.UaaEndpoint()
	if uaaEndpoint == "" {
		return errors.New("UAA endpoint missing from config file")
	}

	url := fmt.Sprintf("%s/Users/%s/password", uaaEndpoint, repo.config.UserGuid())
	body := fmt.Sprintf(`{"password":"%s","oldPassword":"%s"}`, new, old)

	return repo.gateway.UpdateResource(url, strings.NewReader(body))
}
