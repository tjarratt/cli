package app_test

import (
	"github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/app"
	"github.com/tjarratt/cli/cf/command_factory"
	"github.com/tjarratt/cli/cf/io_helpers"
	"github.com/tjarratt/cli/cf/manifest"
	"github.com/tjarratt/cli/cf/net"
	testconfig "github.com/tjarratt/cli/testhelpers/configuration"
	testterm "github.com/tjarratt/cli/testhelpers/terminal"
	"github.com/codegangsta/cli"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
	"time"
)

var _ = Describe("Help", func() {
	It("shows help for all commands", func() {
		commandFactory := createCommandFactory()

		dummyTemplate := `
{{range .Commands}}{{range .CommandSubGroups}}{{range .}}
{{.Name}}
{{end}}{{end}}{{end}}
`
		output := io_helpers.CaptureOutput(func() {
			app.ShowAppHelp(dummyTemplate, createApp(commandFactory))
		})

		for _, metadata := range commandFactory.CommandMetadatas() {
			Expect(commandInOutput(metadata.Name, output)).To(BeTrue(), metadata.Name+" not in help")
		}
	})
})

func createCommandFactory() command_factory.Factory {
	fakeUI := &testterm.FakeUI{}
	configRepo := testconfig.NewRepository()
	manifestRepo := manifest.NewManifestDiskRepository()
	apiRepoLocator := api.NewRepositoryLocator(configRepo, map[string]net.Gateway{
		"auth":             net.NewUAAGateway(configRepo),
		"cloud-controller": net.NewCloudControllerGateway(configRepo, time.Now),
		"uaa":              net.NewUAAGateway(configRepo),
	})

	return command_factory.NewFactory(fakeUI, configRepo, manifestRepo, apiRepoLocator)
}

func createApp(commandFactory command_factory.Factory) *cli.App {
	new_app := cli.NewApp()
	new_app.Commands = []cli.Command{}
	for _, metadata := range commandFactory.CommandMetadatas() {
		new_app.Commands = append(new_app.Commands, cli.Command{Name: metadata.Name})
	}

	return new_app
}

func commandInOutput(cmdName string, output []string) bool {
	for _, line := range output {
		if strings.TrimSpace(line) == strings.TrimSpace(cmdName) {
			return true
		}
	}
	return false
}
