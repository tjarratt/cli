package route

import (
	"cf/api"
	"cf/configuration"
	cferrors "cf/errors"
	"cf/models"
	"cf/requirements"
	"cf/terminal"
	"errors"
	"github.com/tjarratt/cli"
)

type RouteCreator interface {
	CreateRoute(hostName string, domain models.DomainFields, space models.SpaceFields) (route models.Route, apiErr cferrors.Error)
}

type CreateRoute struct {
	ui        terminal.UI
	config    configuration.Reader
	routeRepo api.RouteRepository
	spaceReq  requirements.SpaceRequirement
	domainReq requirements.DomainRequirement
}

func NewCreateRoute(ui terminal.UI, config configuration.Reader, routeRepo api.RouteRepository) (cmd *CreateRoute) {
	cmd = new(CreateRoute)
	cmd.ui = ui
	cmd.config = config
	cmd.routeRepo = routeRepo
	return
}

func (cmd *CreateRoute) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {

	if len(c.Args()) != 2 {
		err = errors.New("Incorrect Usage")
		cmd.ui.FailWithUsage(c, "create-route")
		return
	}

	spaceName := c.Args()[0]
	domainName := c.Args()[1]

	cmd.spaceReq = reqFactory.NewSpaceRequirement(spaceName)
	cmd.domainReq = reqFactory.NewDomainRequirement(domainName)

	reqs = []requirements.Requirement{
		reqFactory.NewLoginRequirement(),
		reqFactory.NewTargetedOrgRequirement(),
		cmd.spaceReq,
		cmd.domainReq,
	}
	return
}

func (cmd *CreateRoute) Run(c *cli.Context) {
	hostName := c.String("n")
	space := cmd.spaceReq.GetSpace()
	domain := cmd.domainReq.GetDomain()

	_, apiErr := cmd.CreateRoute(hostName, domain, space.SpaceFields)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}
}

func (cmd *CreateRoute) CreateRoute(hostName string, domain models.DomainFields, space models.SpaceFields) (route models.Route, apiErr cferrors.Error) {
	cmd.ui.Say("Creating route %s for org %s / space %s as %s...",
		terminal.EntityNameColor(domain.UrlForHost(hostName)),
		terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
		terminal.EntityNameColor(space.Name),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	route, apiErr = cmd.routeRepo.CreateInSpace(hostName, domain.Guid, space.Guid)
	if apiErr != nil {
		var findApiResponse cferrors.Error
		route, findApiResponse = cmd.routeRepo.FindByHostAndDomain(hostName, domain.Name)

		if findApiResponse != nil ||
			route.Space.Guid != space.Guid ||
			route.Domain.Guid != domain.Guid {
			return
		}

		apiErr = nil
		cmd.ui.Ok()
		cmd.ui.Warn("Route %s already exists", route.URL())
		return
	}

	cmd.ui.Ok()
	return
}
