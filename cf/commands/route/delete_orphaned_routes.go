package route

import (
	"github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/models"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/tjarratt/cli/cf/terminal"
	cli "github.com/tjarratt/cg_cli"
)

type DeleteOrphanedRoutes struct {
	ui        terminal.UI
	routeRepo api.RouteRepository
	config    configuration.Reader
}

func NewDeleteOrphanedRoutes(ui terminal.UI, config configuration.Reader, routeRepo api.RouteRepository) (cmd DeleteOrphanedRoutes) {
	cmd.ui = ui
	cmd.config = config
	cmd.routeRepo = routeRepo
	return
}

func (cmd DeleteOrphanedRoutes) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "delete-orphaned-routes",
		Description: "Delete all orphaned routes (e.g.: those that are not mapped to an app)",
		Usage:       "CF_NAME delete-orphaned-routes [-f]",
		Flags: []cli.Flag{
			cli.BoolFlag{Name: "f", Usage: "Force deletion without confirmation"},
		},
	}
}

func (cmd DeleteOrphanedRoutes) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	reqs = append(reqs, reqFactory.NewLoginRequirement())
	return
}

func (cmd DeleteOrphanedRoutes) Run(c *cli.Context) {

	force := c.Bool("f")
	if !force {
		response := cmd.ui.Confirm(
			"Really delete orphaned routes?%s",
			terminal.PromptColor(">"),
		)

		if !response {
			return
		}
	}

	cmd.ui.Say("Getting routes as %s ...\n",
		terminal.EntityNameColor(cmd.config.Username()),
	)

	apiErr := cmd.routeRepo.ListRoutes(func(route models.Route) bool {

		if len(route.Apps) == 0 {
			cmd.ui.Say("Deleting route %s...", terminal.EntityNameColor(route.Host+"."+route.Domain.Name))
			apiErr := cmd.routeRepo.Delete(route.Guid)
			if apiErr != nil {
				cmd.ui.Failed(apiErr.Error())
				return false
			}
		}
		return true
	})

	if apiErr != nil {
		cmd.ui.Failed("Failed fetching routes.\n%s", apiErr.Error())
		return
	}
	cmd.ui.Ok()
}
