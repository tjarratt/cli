package commands

import (
	"flag"
	"fmt"
	"github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/app"
	"github.com/tjarratt/cli/cf/command_factory"
	"github.com/tjarratt/cli/cf/command_runner"
	"github.com/tjarratt/cli/cf/manifest"
	"github.com/tjarratt/cli/cf/net"
	testconfig "github.com/tjarratt/cli/testhelpers/configuration"
	testreq "github.com/tjarratt/cli/testhelpers/requirements"
	testterm "github.com/tjarratt/cli/testhelpers/terminal"
	cli "github.com/tjarratt/cg_cli"
	"strings"
	"time"
)

func NewContext(cmdName string, args []string) *cli.Context {
	targetCommand := findCommand(cmdName)

	flagSet := new(flag.FlagSet)
	for i, _ := range targetCommand.Flags {
		targetCommand.Flags[i].Apply(flagSet)
	}

	// move all flag args to the beginning of the list, go requires them all upfront
	firstFlagIndex := -1
	for index, arg := range args {
		if strings.HasPrefix(arg, "-") {
			firstFlagIndex = index
			break
		}
	}
	if firstFlagIndex > 0 {
		args := args[0:firstFlagIndex]
		flags := args[firstFlagIndex:]
		flagSet.Parse(append(flags, args...))
	} else {
		flagSet.Parse(args[0:])
	}

	globalSet := new(flag.FlagSet)

	return cli.NewContext(cli.NewApp(), flagSet, globalSet)
}

func findCommand(cmdName string) (cmd cli.Command) {
	fakeUI := &testterm.FakeUI{}
	configRepo := testconfig.NewRepository()
	manifestRepo := manifest.NewManifestDiskRepository()
	apiRepoLocator := api.NewRepositoryLocator(configRepo, map[string]net.Gateway{
		"auth":             net.NewUAAGateway(configRepo),
		"cloud-controller": net.NewCloudControllerGateway(configRepo, time.Now),
		"uaa":              net.NewUAAGateway(configRepo),
	})

	cmdFactory := command_factory.NewFactory(fakeUI, configRepo, manifestRepo, apiRepoLocator)
	requirementsFactory := &testreq.FakeReqFactory{}
	cmdRunner := command_runner.NewRunner(cmdFactory, requirementsFactory)
	myApp := app.NewApp(cmdRunner, cmdFactory.CommandMetadatas()...)

	for _, cmd := range myApp.Commands {
		if cmd.Name == cmdName {
			return cmd
		}
	}
	panic(fmt.Sprintf("command %s does not exist", cmdName))
	return
}
