package application

import (
	cli "github.com/tjarratt/cg_cli"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/models"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/tjarratt/cli/cf/terminal"
)

type Restart struct {
	ui      terminal.UI
	starter ApplicationStarter
	stopper ApplicationStopper
	appReq  requirements.ApplicationRequirement
}

type ApplicationRestarter interface {
	ApplicationRestart(app models.Application)
}

func NewRestart(ui terminal.UI, starter ApplicationStarter, stopper ApplicationStopper) (cmd *Restart) {
	cmd = new(Restart)
	cmd.ui = ui
	cmd.starter = starter
	cmd.stopper = stopper
	return
}

func (cmd *Restart) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "restart",
		ShortName:   "rs",
		Description: "Restart an app",
		Usage:       "CF_NAME restart APP",
	}
}

func (cmd *Restart) GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) == 0 {
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

func (cmd *Restart) Run(c *cli.Context) {
	app := cmd.appReq.GetApplication()
	cmd.ApplicationRestart(app)
}

func (cmd *Restart) ApplicationRestart(app models.Application) {
	stoppedApp, err := cmd.stopper.ApplicationStop(app)
	if err != nil {
		cmd.ui.Failed(err.Error())
		return
	}

	cmd.ui.Say("")

	_, err = cmd.starter.ApplicationStart(stoppedApp)
	if err != nil {
		cmd.ui.Failed(err.Error())
		return
	}
}
