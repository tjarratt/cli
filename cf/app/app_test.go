package app_test

import (
	"bytes"
	"strings"
	"time"

	cli "github.com/tjarratt/cg_cli"
	"github.com/tjarratt/cli/cf"
	"github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/command_factory"
	"github.com/tjarratt/cli/cf/command_metadata"
	"github.com/tjarratt/cli/cf/net"
	"github.com/tjarratt/cli/cf/trace"
	testconfig "github.com/tjarratt/cli/testhelpers/configuration"
	testmanifest "github.com/tjarratt/cli/testhelpers/manifest"
	testterm "github.com/tjarratt/cli/testhelpers/terminal"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/tjarratt/cli/cf/app"
	. "github.com/tjarratt/cli/testhelpers/matchers"
)

var expectedCommandNames = []string{
	"api", "app", "apps", "auth", "bind-service", "buildpacks", "create-buildpack",
	"create-domain", "create-org", "create-route", "create-service", "create-service-auth-token",
	"create-service-broker", "create-space", "create-user", "create-user-provided-service", "curl",
	"delete", "delete-buildpack", "delete-domain", "delete-shared-domain", "delete-org", "delete-route",
	"delete-service", "delete-service-auth-token", "delete-service-broker", "delete-space", "delete-user",
	"domains", "env", "events", "files", "login", "logout", "logs", "marketplace", "map-route", "org",
	"org-users", "orgs", "passwd", "purge-service-offering", "push", "quotas", "rename", "rename-org",
	"rename-service", "rename-service-broker", "rename-space", "restart", "routes", "scale",
	"service", "service-auth-tokens", "service-brokers", "services", "set-env", "set-org-role",
	"set-space-role", "create-shared-domain", "space", "space-users", "spaces", "stacks", "start", "stop",
	"target", "unbind-service", "unmap-route", "unset-env", "unset-org-role", "unset-space-role",
	"update-buildpack", "update-service-broker", "update-service-auth-token", "update-user-provided-service",
	"quotas", "create-quota", "delete-quota", "quota", "set-quota",
}

var _ = Describe("App", func() {
	var (
		cmdRunner  *FakeRunner
		cmdFactory command_factory.Factory
		app        *cli.App
	)

	BeforeEach(func() {
		ui := &testterm.FakeUI{}
		config := testconfig.NewRepository()
		manifestRepo := &testmanifest.FakeManifestRepository{}

		repoLocator := api.NewRepositoryLocator(config, map[string]net.Gateway{
			"auth":             net.NewUAAGateway(config),
			"cloud-controller": net.NewCloudControllerGateway(config, time.Now),
			"uaa":              net.NewUAAGateway(config),
		})

		cmdFactory = command_factory.NewFactory(ui, config, manifestRepo, repoLocator)

		metadatas := []command_metadata.CommandMetadata{}
		for _, cmdName := range expectedCommandNames {
			metadatas = append(metadatas, command_metadata.CommandMetadata{Name: cmdName})
		}

		cmdRunner = &FakeRunner{cmdFactory: cmdFactory}

	})

	JustBeforeEach(func() {
		app = NewApp(cmdRunner, cmdFactory.CommandMetadatas()...)
	})

	Describe("tracefile integration", func() {
		var output *bytes.Buffer

		BeforeEach(func() {
			output = bytes.NewBuffer(make([]byte, 1024))
			trace.SetStdout(output)
			trace.EnableTrace()
		})

		It("prints its version during its constructor", func() {
			Expect(strings.Split(output.String(), "\n")).To(ContainSubstrings(
				[]string{"VERSION:"},
				[]string{cf.Version},
			))
		})
	})

	It("#NewApp", func() {
		for _, cmdName := range expectedCommandNames {
			app.Run([]string{"hey-look-ma-no-hands", cmdName})
			Expect(cmdRunner.cmdName).To(Equal(cmdName))
		}
	})

	Describe("autocorrect", func() {
		It("will change halp to help", func() {
			app.Run([]string{"scott-gcf", "halp", "me"})
			Expect(cmdRunner.cmdName).To(Equal("help"))
		})

		It("passes the rest of the args to the correct command", func() {
			app.Run([]string{"scott-gcf", "halp", "one", "two", "buckle mah shoe"})
			Expect(cmdRunner.cmdArgs).To(Equal([]string{"one", "two", "buckle mah shoe"}))
		})

		It("will change jalp to help", func() {
			app.Run([]string{"scott-gcf", "jalp", "help"})
			Expect(cmdRunner.cmdName).To(Equal("help"))
		})
	})
})

type FakeRunner struct {
	cmdFactory command_factory.Factory
	cmdName    string
	cmdArgs    []string
}

func (runner *FakeRunner) RunCmdByName(cmdName string, c *cli.Context) (err error) {
	_, err = runner.cmdFactory.GetByCmdName(cmdName)
	if err != nil {
		GinkgoT().Fatal("Error instantiating command with name", cmdName)
		return
	}
	runner.cmdName = cmdName
	runner.cmdArgs = c.Args()
	return
}
