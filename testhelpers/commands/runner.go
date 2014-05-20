package commands

import (
	"github.com/tjarratt/cli/cf/command"
	testreq "github.com/tjarratt/cli/testhelpers/requirements"
	testterm "github.com/tjarratt/cli/testhelpers/terminal"
	cli "github.com/tjarratt/cg_cli"
)

var CommandDidPassRequirements bool

func RunCommand(cmd command.Command, ctxt *cli.Context, requirementsFactory *testreq.FakeReqFactory) (passedRequirements bool) {
	defer func() {
		errMsg := recover()

		if errMsg != nil && errMsg != testterm.FailedWasCalled {
			panic(errMsg)
		}
	}()

	CommandDidPassRequirements = false

	requirements, err := cmd.GetRequirements(requirementsFactory, ctxt)
	if err != nil {
		return
	}

	for _, requirement := range requirements {
		success := requirement.Execute()
		if !success {
			return
		}
	}

	passedRequirements = true
	CommandDidPassRequirements = true
	cmd.Run(ctxt)

	return
}
