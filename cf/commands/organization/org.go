package organization

import (
	"fmt"
	cli "github.com/tjarratt/cg_cli"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/formatters"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/tjarratt/cli/cf/terminal"
	"strings"
)

type ShowOrg struct {
	ui     terminal.UI
	config configuration.Reader
	orgReq requirements.OrganizationRequirement
}

func NewShowOrg(ui terminal.UI, config configuration.Reader) (cmd *ShowOrg) {
	cmd = new(ShowOrg)
	cmd.ui = ui
	cmd.config = config
	return
}

func (cmd *ShowOrg) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "org",
		Description: "Show org info",
		Usage:       "CF_NAME org ORG",
	}
}

func (cmd *ShowOrg) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 1 {
		cmd.ui.FailWithUsage(c)
	}

	cmd.orgReq = requirementsFactory.NewOrganizationRequirement(c.Args()[0])
	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		cmd.orgReq,
	}

	return
}

func (cmd *ShowOrg) Run(c *cli.Context) {
	org := cmd.orgReq.GetOrganization()
	cmd.ui.Say("Getting info for org %s as %s...",
		terminal.EntityNameColor(org.Name),
		terminal.EntityNameColor(cmd.config.Username()),
	)
	cmd.ui.Ok()
	cmd.ui.Say("\n%s:", terminal.EntityNameColor(org.Name))

	domains := []string{}
	for _, domain := range org.Domains {
		domains = append(domains, domain.Name)
	}

	spaces := []string{}
	for _, space := range org.Spaces {
		spaces = append(spaces, space.Name)
	}

	quota := org.QuotaDefinition
	orgQuota := fmt.Sprintf("%s (%dM memory limit, %d routes, %d services, paid services %s)",
		quota.Name, quota.MemoryLimit, quota.RoutesLimit, quota.ServicesLimit, formatters.Allowed(quota.NonBasicServicesAllowed))

	cmd.ui.Say("  domains: %s", terminal.EntityNameColor(strings.Join(domains, ", ")))
	cmd.ui.Say("  quota:   %s", terminal.EntityNameColor(orgQuota))
	cmd.ui.Say("  spaces:  %s", terminal.EntityNameColor(strings.Join(spaces, ", ")))
}
