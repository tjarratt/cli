package service_test

import (
	. "github.com/tjarratt/cli/cf/commands/service"
	"github.com/tjarratt/cli/cf/configuration"
	testapi "github.com/tjarratt/cli/testhelpers/api"
	testcmd "github.com/tjarratt/cli/testhelpers/commands"
	testconfig "github.com/tjarratt/cli/testhelpers/configuration"
	testreq "github.com/tjarratt/cli/testhelpers/requirements"
	testterm "github.com/tjarratt/cli/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/tjarratt/cli/testhelpers/matchers"
)

var _ = Describe("create-user-provided-service command", func() {
	var (
		ui                  *testterm.FakeUI
		config              configuration.ReadWriter
		repo                *testapi.FakeUserProvidedServiceInstanceRepo
		requirementsFactory *testreq.FakeReqFactory
		cmd                 CreateUserProvidedService
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		config = testconfig.NewRepositoryWithDefaults()
		repo = &testapi.FakeUserProvidedServiceInstanceRepo{}
		requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: true}
		cmd = NewCreateUserProvidedService(ui, config, repo)
	})

	Describe("login requirements", func() {
		It("fails if the user is not logged in", func() {
			requirementsFactory.LoginSuccess = false
			ctxt := testcmd.NewContext("create-user-provided-service", []string{"my-service"})
			testcmd.RunCommand(cmd, ctxt, requirementsFactory)
			Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
		})
	})

	It("creates a new user provided service given just a name", func() {
		args := []string{"my-custom-service"}
		ctxt := testcmd.NewContext("create-user-provided-service", args)
		testcmd.RunCommand(cmd, ctxt, requirementsFactory)
		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Creating user provided service"},
			[]string{"OK"},
		))
	})

	It("accepts service parameters interactively", func() {
		ui.Inputs = []string{"foo value", "bar value", "baz value"}
		ctxt := testcmd.NewContext("create-user-provided-service", []string{"-p", `"foo, bar, baz"`, "my-custom-service"})
		testcmd.RunCommand(cmd, ctxt, requirementsFactory)

		Expect(ui.Prompts).To(ContainSubstrings(
			[]string{"foo"},
			[]string{"bar"},
			[]string{"baz"},
		))

		Expect(repo.CreateName).To(Equal("my-custom-service"))
		Expect(repo.CreateParams).To(Equal(map[string]string{
			"foo": "foo value",
			"bar": "bar value",
			"baz": "baz value",
		}))

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Creating user provided service", "my-custom-service", "my-org", "my-space", "my-user"},
			[]string{"OK"},
		))
	})

	It("accepts service parameters as JSON without prompting", func() {
		args := []string{"-p", `{"foo": "foo value", "bar": "bar value", "baz": "baz value"}`, "my-custom-service"}
		ctxt := testcmd.NewContext("create-user-provided-service", args)
		testcmd.RunCommand(cmd, ctxt, requirementsFactory)

		Expect(ui.Prompts).To(BeEmpty())
		Expect(repo.CreateName).To(Equal("my-custom-service"))
		Expect(repo.CreateParams).To(Equal(map[string]string{
			"foo": "foo value",
			"bar": "bar value",
			"baz": "baz value",
		}))

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Creating user provided service"},
			[]string{"OK"},
		))
	})

	It("creates a user provided service with a syslog drain url", func() {
		args := []string{"-l", "syslog://example.com", "-p", `{"foo": "foo value", "bar": "bar value", "baz": "baz value"}`, "my-custom-service"}
		ctxt := testcmd.NewContext("create-user-provided-service", args)
		testcmd.RunCommand(cmd, ctxt, requirementsFactory)

		Expect(repo.CreateDrainUrl).To(Equal("syslog://example.com"))
		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Creating user provided service"},
			[]string{"OK"},
		))
	})
})
