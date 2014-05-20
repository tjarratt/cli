package command_metadata

import cli "github.com/tjarratt/cg_cli"

type CommandMetadata struct {
	Name            string
	ShortName       string
	Usage           string
	Description     string
	Flags           []cli.Flag
	SkipFlagParsing bool
}
