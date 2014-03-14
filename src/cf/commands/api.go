package commands

import (
	"cf"
	"cf/api"
	"cf/configuration"
	"cf/errors"
	"cf/requirements"
	"cf/terminal"
	"fmt"
	"github.com/tjarratt/cli"
	"strings"
)

type Api struct {
	ui           terminal.UI
	endpointRepo api.EndpointRepository
	config       configuration.ReadWriter
}

func NewApi(ui terminal.UI, config configuration.ReadWriter, endpointRepo api.EndpointRepository) (cmd Api) {
	cmd.ui = ui
	cmd.config = config
	cmd.endpointRepo = endpointRepo
	return
}

func (cmd Api) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	return
}

func (cmd Api) Run(c *cli.Context) {
	if len(c.Args()) == 0 {
		if cmd.config.ApiEndpoint() == "" {
			cmd.ui.Say(fmt.Sprintf("No api endpoint set. Use '%s' to set an endpoint", terminal.CommandColor(cf.Name()+" api")))
		} else {
			cmd.ui.Say(
				"API endpoint: %s (API version: %s)",
				terminal.EntityNameColor(cmd.config.ApiEndpoint()),
				terminal.EntityNameColor(cmd.config.ApiVersion()),
			)
		}
		return
	}

	endpoint := c.Args()[0]

	cmd.ui.Say("Setting api endpoint to %s...", terminal.EntityNameColor(endpoint))
	cmd.setApiEndpoint(endpoint, c.Bool("skip-ssl-validation"))
	cmd.ui.Ok()

	cmd.ui.Say("")
	cmd.ui.ShowConfiguration(cmd.config)
}

func (cmd Api) setApiEndpoint(endpoint string, skipSSL bool) {
	if strings.HasSuffix(endpoint, "/") {
		endpoint = strings.TrimSuffix(endpoint, "/")
	}

	cmd.config.SetSSLDisabled(skipSSL)
	endpoint, err := cmd.endpointRepo.UpdateEndpoint(endpoint)

	if err != nil {
		cmd.config.SetApiEndpoint("")
		cmd.config.SetSSLDisabled(false)

		switch typedErr := err.(type) {
		case *errors.InvalidSSLCert:
			cfApiCommand := terminal.CommandColor(fmt.Sprintf("%s api --skip-ssl-validation", cf.Name()))
			tipMessage := fmt.Sprintf("TIP: Use '%s' to continue with an insecure API endpoint", cfApiCommand)
			cmd.ui.Failed("Invalid SSL Cert for %s\n%s", typedErr.URL, tipMessage)
		default:
			cmd.ui.Failed(typedErr.Error())
		}
	}

	if !strings.HasPrefix(endpoint, "https://") {
		cmd.ui.Say(terminal.WarningColor("Warning: Insecure http API endpoint detected: secure https API endpoints are recommended\n"))
	}
}
