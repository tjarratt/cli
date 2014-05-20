package domain

import (
	"github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/tjarratt/cli/cf/terminal"
	cli "github.com/tjarratt/cg_cli"
)

type CreateSharedDomain struct {
	ui         terminal.UI
	config     configuration.Reader
	domainRepo api.DomainRepository
	orgReq     requirements.OrganizationRequirement
}

func NewCreateSharedDomain(ui terminal.UI, config configuration.Reader, domainRepo api.DomainRepository) (cmd *CreateSharedDomain) {
	cmd = new(CreateSharedDomain)
	cmd.ui = ui
	cmd.config = config
	cmd.domainRepo = domainRepo
	return
}

func (cmd *CreateSharedDomain) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "create-shared-domain",
		Description: "Create a domain that can be used by all orgs (admin-only)",
		Usage:       "CF_NAME create-shared-domain DOMAIN",
	}
}

func (cmd *CreateSharedDomain) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 1 {
		cmd.ui.FailWithUsage(c)
	}

	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
	}
	return
}

func (cmd *CreateSharedDomain) Run(c *cli.Context) {
	domainName := c.Args()[0]

	cmd.ui.Say("Creating shared domain %s as %s...",
		terminal.EntityNameColor(domainName),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	apiErr := cmd.domainRepo.CreateSharedDomain(domainName)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
}
