package securitygroup

import (
	"github.com/cloudfoundry/cli/cf/api/security_groups/defaults/staging"
	"github.com/cloudfoundry/cli/cf/command_metadata"
	"github.com/cloudfoundry/cli/cf/configuration"
	"github.com/cloudfoundry/cli/cf/requirements"
	"github.com/cloudfoundry/cli/cf/terminal"
	"github.com/codegangsta/cli"
)

type listStagingSecurityGroups struct {
	ui                       terminal.UI
	stagingSecurityGroupRepo staging.StagingSecurityGroupsRepo
	configRepo               configuration.Reader
}

func NewListStagingSecurityGroups(ui terminal.UI, configRepo configuration.Reader, stagingSecurityGroupRepo staging.StagingSecurityGroupsRepo) listStagingSecurityGroups {
	return listStagingSecurityGroups{
		ui:                       ui,
		configRepo:               configRepo,
		stagingSecurityGroupRepo: stagingSecurityGroupRepo,
	}
}

func (cmd listStagingSecurityGroups) Metadata() command_metadata.CommandMetadata {
	return command_metadata.CommandMetadata{
		Name:        "staging-security-groups",
		Description: "List security groups in the staging set for applications",
		Usage:       "CF_NAME staging-security-groups",
	}
}

func (cmd listStagingSecurityGroups) GetRequirements(requirementsFactory requirements.Factory, context *cli.Context) ([]requirements.Requirement, error) {
	if len(context.Args()) != 0 {
		cmd.ui.FailWithUsage(context)
	}

	requirements := []requirements.Requirement{requirementsFactory.NewLoginRequirement()}
	return requirements, nil
}

func (cmd listStagingSecurityGroups) Run(context *cli.Context) {
	cmd.ui.Say("Acquiring staging security group as %s",
		terminal.EntityNameColor(cmd.configRepo.Username()))

	SecurityGroupsFields, err := cmd.stagingSecurityGroupRepo.List()
	if err != nil {
		cmd.ui.Failed(err.Error())
	}

	cmd.ui.Ok()
	cmd.ui.Say("")

	if len(SecurityGroupsFields) > 0 {
		for _, value := range SecurityGroupsFields {
			cmd.ui.Say(value.Name)
		}
	} else {
		cmd.ui.Say("No staging security group set")
	}
}
