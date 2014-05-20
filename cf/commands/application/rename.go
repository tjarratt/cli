package application

import (
	"github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/models"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/tjarratt/cli/cf/terminal"
	cli "github.com/tjarratt/cg_cli"
)

type RenameApp struct {
	ui      terminal.UI
	config  configuration.Reader
	appRepo api.ApplicationRepository
	appReq  requirements.ApplicationRequirement
}

func NewRenameApp(ui terminal.UI, config configuration.Reader, appRepo api.ApplicationRepository) (cmd *RenameApp) {
	cmd = new(RenameApp)
	cmd.ui = ui
	cmd.config = config
	cmd.appRepo = appRepo
	return
}

func (cmd *RenameApp) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "rename",
		Description: "Rename an app",
		Usage:       "CF_NAME rename APP_NAME NEW_APP_NAME",
	}
}

func (cmd *RenameApp) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 2 {
		cmd.ui.FailWithUsage(c)
	}

	cmd.appReq = requirementsFactory.NewApplicationRequirement(c.Args()[0])
	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		cmd.appReq,
	}
	return
}

func (cmd *RenameApp) Run(c *cli.Context) {
	app := cmd.appReq.GetApplication()
	newName := c.Args()[1]

	cmd.ui.Say("Renaming app %s to %s in org %s / space %s as %s...",
		terminal.EntityNameColor(app.Name),
		terminal.EntityNameColor(newName),
		terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
		terminal.EntityNameColor(cmd.config.SpaceFields().Name),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	params := models.AppParams{Name: &newName}

	_, apiErr := cmd.appRepo.Update(app.Guid, params)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}
	cmd.ui.Ok()
}
