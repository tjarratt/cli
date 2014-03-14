package terminal

import (
	"cf/configuration"
	term "cf/terminal"
	"fmt"
	"github.com/tjarratt/cli"
	"strings"
	"time"
)

const FailedWasCalled = "FailedWasCalled"

type FakeUI struct {
	Outputs                    []string
	Prompts                    []string
	PasswordPrompts            []string
	Inputs                     []string
	FailedWithUsage            bool
	FailedWithUsageCommandName string
	ShowConfigurationCalled    bool
}

func (ui *FakeUI) PrintPaginator(rows []string, err error) {
	if err != nil {
		ui.Failed(err.Error())
		return
	}

	for _, row := range rows {
		ui.Say(row)
	}
}

func (ui *FakeUI) Say(message string, args ...interface{}) {
	message = fmt.Sprintf(message, args...)
	ui.Outputs = append(ui.Outputs, strings.Split(message, "\n")...)
	return
}

func (ui *FakeUI) Warn(message string, args ...interface{}) {
	ui.Say(message, args...)
	return
}

func (ui *FakeUI) Ask(prompt string, args ...interface{}) (answer string) {
	ui.Prompts = append(ui.Prompts, fmt.Sprintf(prompt, args...))
	if len(ui.Inputs) == 0 {
		panic("No input provided to Fake UI for prompt: " + fmt.Sprintf(prompt, args...))
	}

	answer = ui.Inputs[0]
	ui.Inputs = ui.Inputs[1:]
	return
}

func (ui *FakeUI) Confirm(prompt string, args ...interface{}) bool {
	response := ui.Ask(prompt, args...)
	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	}
	return false
}

func (ui *FakeUI) AskForPassword(prompt string, args ...interface{}) (answer string) {
	ui.PasswordPrompts = append(ui.PasswordPrompts, fmt.Sprintf(prompt, args...))
	if len(ui.Inputs) == 0 {
		panic("No input provided to Fake UI for prompt: " + fmt.Sprintf(prompt, args...))
	}

	answer = ui.Inputs[0]
	ui.Inputs = ui.Inputs[1:]
	return
}

func (ui *FakeUI) Ok() {
	ui.Say("OK")
}

func (ui *FakeUI) Failed(message string, args ...interface{}) {
	ui.Say("FAILED")
	ui.Say(message, args...)
	panic(FailedWasCalled)
	return
}

func (ui *FakeUI) ConfigFailure(err error) {
	ui.Failed("Error loading config file.\n%s", err.Error())
}

func (ui *FakeUI) FailWithUsage(ctxt *cli.Context, cmdName string) {
	ui.FailedWithUsage = true
	ui.FailedWithUsageCommandName = cmdName
	ui.Failed("Incorrect Usage.")
}

func (ui *FakeUI) DumpOutputs() string {
	return "****************************\n" + strings.Join(ui.Outputs, "\n")
}

func (ui *FakeUI) ClearOutputs() {
	ui.Outputs = []string{}
}

func (ui *FakeUI) ShowConfiguration(config configuration.Reader) {
	ui.ShowConfigurationCalled = true
}

func (ui FakeUI) LoadingIndication() {
}

func (c FakeUI) Wait(duration time.Duration) {
	time.Sleep(duration)
}

func (ui *FakeUI) DisplayTable(table [][]string) {

	for _, line := range table {
		output := ""
		for _, value := range line {
			output = output + value + "  "
		}
		ui.Say("%s", output)
	}
}

func (ui *FakeUI) Table(headers []string) term.Table {
	return term.NewTable(ui, headers)
}
