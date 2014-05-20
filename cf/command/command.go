package command

import (
	cli "github.com/tjarratt/cg_cli"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/requirements"
)

type Command interface {
	Metadata() command_metadata.CommandMetadata
	GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error)
	Run(c *cli.Context)
}
