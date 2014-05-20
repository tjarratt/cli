package command

import (
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/requirements"
	"github.com/codegangsta/cli"
)

type Command interface {
	Metadata() command_metadata.CommandMetadata
	GetRequirements(requirementsFactory requirements.Factory, c *cli.Context) (reqs []requirements.Requirement, err error)
	Run(c *cli.Context)
}
