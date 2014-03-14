package application

import (
	"cf/api"
	"cf/configuration"
	"cf/requirements"
	"cf/terminal"
	"errors"
	"github.com/tjarratt/cli"
)

type Files struct {
	ui           terminal.UI
	config       configuration.Reader
	appFilesRepo api.AppFilesRepository
	appReq       requirements.ApplicationRequirement
}

func NewFiles(ui terminal.UI, config configuration.Reader, appFilesRepo api.AppFilesRepository) (cmd *Files) {
	cmd = new(Files)
	cmd.ui = ui
	cmd.config = config
	cmd.appFilesRepo = appFilesRepo
	return
}

func (cmd *Files) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) < 1 {
		err = errors.New("Incorrect Usage")
		cmd.ui.FailWithUsage(c, "files")
		return
	}

	cmd.appReq = reqFactory.NewApplicationRequirement(c.Args()[0])

	reqs = []requirements.Requirement{
		reqFactory.NewLoginRequirement(),
		reqFactory.NewTargetedSpaceRequirement(),
		cmd.appReq,
	}
	return
}

func (cmd *Files) Run(c *cli.Context) {
	app := cmd.appReq.GetApplication()

	cmd.ui.Say("Getting files for app %s in org %s / space %s as %s...",
		terminal.EntityNameColor(app.Name),
		terminal.EntityNameColor(cmd.config.OrganizationFields().Name),
		terminal.EntityNameColor(cmd.config.SpaceFields().Name),
		terminal.EntityNameColor(cmd.config.Username()),
	)

	path := "/"
	if len(c.Args()) > 1 {
		path = c.Args()[1]
	}

	list, apiErr := cmd.appFilesRepo.ListFiles(app.Guid, path)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
	cmd.ui.Say("")
	cmd.ui.Say("%s", list)
}
