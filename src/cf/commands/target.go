package commands

import (
	"cf/api"
	"cf/configuration"
	"cf/errors"
	"cf/models"
	"cf/requirements"
	"cf/terminal"
	"github.com/tjarratt/cli"
)

type Target struct {
	ui        terminal.UI
	config    configuration.ReadWriter
	orgRepo   api.OrganizationRepository
	spaceRepo api.SpaceRepository
}

func NewTarget(ui terminal.UI,
	config configuration.ReadWriter,
	orgRepo api.OrganizationRepository,
	spaceRepo api.SpaceRepository) (cmd Target) {

	cmd.ui = ui
	cmd.config = config
	cmd.orgRepo = orgRepo
	cmd.spaceRepo = spaceRepo

	return
}

func (cmd Target) GetRequirements(reqFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error) {
	if len(c.Args()) != 0 {
		err = errors.New("incorrect usage")
		cmd.ui.FailWithUsage(c, "target")
		return
	}

	if c.String("o") != "" || c.String("s") != "" {
		reqs = append(reqs, reqFactory.NewLoginRequirement())
	}

	return
}

func (cmd Target) Run(c *cli.Context) {
	orgName := c.String("o")
	spaceName := c.String("s")

	if orgName != "" {
		err := cmd.setOrganization(orgName)
		if err != nil {
			cmd.ui.Failed(err.Error())
		}
	}

	if spaceName != "" {
		err := cmd.setSpace(spaceName)
		if err != nil {
			cmd.ui.Failed(err.Error())
		}
	}

	cmd.ui.ShowConfiguration(cmd.config)
	return
}

func (cmd Target) setOrganization(orgName string) error {
	// setting an org necessarily invalidates any space you had previously targeted
	cmd.config.SetOrganizationFields(models.OrganizationFields{})
	cmd.config.SetSpaceFields(models.SpaceFields{})

	org, apiErr := cmd.orgRepo.FindByName(orgName)
	if apiErr != nil {
		return errors.NewErrorWithMessage("Could not target org.\n%s", apiErr.Error())
	}

	cmd.config.SetOrganizationFields(org.OrganizationFields)
	return nil
}

func (cmd Target) setSpace(spaceName string) error {
	cmd.config.SetSpaceFields(models.SpaceFields{})

	if !cmd.config.HasOrganization() {
		return errors.New("An org must be targeted before targeting a space")
	}

	space, apiErr := cmd.spaceRepo.FindByName(spaceName)
	if apiErr != nil {
		return errors.NewErrorWithMessage("Unable to access space %s.\n%s", spaceName, apiErr.Error())
	}

	cmd.config.SetSpaceFields(space.SpaceFields)
	return nil
}
