package application

import (
	cli "github.com/tjarratt/cg_cli"
	"github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/tjarratt/cli/cf/terminal"
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

func (cmd *Files) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "files",
		ShortName:   "f",
		Description: "Print out a list of files in a directory or the contents of a specific file",
		Usage:       "CF_NAME files APP [PATH]",
	}
}

func (cmd *Files) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) < 1 {
		cmd.ui.FailWithUsage(c)
	}

	cmd.appReq = requirementsFactory.NewApplicationRequirement(c.Args()[0])

	reqs = []requirements.Requirement{
		requirementsFactory.NewLoginRequirement(),
		requirementsFactory.NewTargetedSpaceRequirement(),
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
