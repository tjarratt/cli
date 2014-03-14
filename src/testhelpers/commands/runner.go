package commands

import (
	"cf/commands"
	"github.com/tjarratt/cli"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
)

var CommandDidPassRequirements bool

func RunCommand(cmd commands.Command, ctxt *cli.Context, reqFactory *testreq.FakeReqFactory) {
	defer func() {
		errMsg := recover()

		if errMsg != nil && errMsg != testterm.FailedWasCalled {
			panic(errMsg)
		}
	}()

	CommandDidPassRequirements = false

	requirements, err := cmd.GetRequirements(reqFactory, ctxt)
	if err != nil {
		return
	}

	for _, requirement := range requirements {
		success := requirement.Execute()
		if !success {
			return
		}
	}

	CommandDidPassRequirements = true
	cmd.Run(ctxt)

	return
}
