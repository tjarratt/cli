package buildpack

import (
	"cf"
	"cf/api"
	"cf/errors"
	"cf/models"
	"cf/requirements"
	"cf/terminal"
	"github.com/tjarratt/cli"
	"strconv"
)

type CreateBuildpack struct {
	ui                terminal.UI
	buildpackRepo     api.BuildpackRepository
	buildpackBitsRepo api.BuildpackBitsRepository
}

func NewCreateBuildpack(ui terminal.UI, buildpackRepo api.BuildpackRepository, buildpackBitsRepo api.BuildpackBitsRepository) (cmd CreateBuildpack) {
	cmd.ui = ui
	cmd.buildpackRepo = buildpackRepo
	cmd.buildpackBitsRepo = buildpackBitsRepo
	return
}

func (cmd CreateBuildpack) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	reqs = []requirements.Requirement{
		reqFactory.NewLoginRequirement(),
	}
	return
}

func (cmd CreateBuildpack) Run(c *cli.Context) {
	if len(c.Args()) != 3 {
		cmd.ui.FailWithUsage(c, "create-buildpack")
		return
	}

	buildpackName := c.Args()[0]

	cmd.ui.Say("Creating buildpack %s...", terminal.EntityNameColor(buildpackName))

	buildpack, apiErr := cmd.createBuildpack(buildpackName, c)
	if apiErr != nil {
		if apiErr.ErrorCode() == cf.BUILDPACK_EXISTS {
			cmd.ui.Ok()
			cmd.ui.Warn("Buildpack %s already exists", buildpackName)
			cmd.ui.Say("TIP: use '%s' to update this buildpack", terminal.CommandColor(cf.Name()+" update-buildpack"))
		} else {
			cmd.ui.Failed(apiErr.Error())
		}
		return
	}
	cmd.ui.Ok()
	cmd.ui.Say("")

	cmd.ui.Say("Uploading buildpack %s...", terminal.EntityNameColor(buildpackName))

	dir := c.Args()[1]

	apiErr = cmd.buildpackBitsRepo.UploadBuildpack(buildpack, dir)
	if apiErr != nil {
		cmd.ui.Failed(apiErr.Error())
		return
	}

	cmd.ui.Ok()
}

func (cmd CreateBuildpack) createBuildpack(buildpackName string, c *cli.Context) (buildpack models.Buildpack, apiErr errors.Error) {
	position, err := strconv.Atoi(c.Args()[2])
	if err != nil {
		apiErr = errors.NewErrorWithMessage("Invalid position. %s", err.Error())
		return
	}

	enabled := c.Bool("enable")
	disabled := c.Bool("disable")
	if enabled && disabled {
		apiErr = errors.NewErrorWithMessage("Cannot specify both enabled and disabled.")
		return
	}

	var enableOption *bool = nil
	if enabled {
		enableOption = &enabled
	}
	if disabled {
		disabled = false
		enableOption = &disabled
	}

	buildpack, apiErr = cmd.buildpackRepo.Create(buildpackName, &position, enableOption, nil)

	return
}
