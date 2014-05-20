package buildpack_test

import (
	"github.com/tjarratt/cli/cf"
	testapi "github.com/tjarratt/cli/testhelpers/api"
	testcmd "github.com/tjarratt/cli/testhelpers/commands"
	testreq "github.com/tjarratt/cli/testhelpers/requirements"
	testterm "github.com/tjarratt/cli/testhelpers/terminal"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/tjarratt/cli/cf/commands/buildpack"
	. "github.com/tjarratt/cli/testhelpers/matchers"
)

var _ = Describe("create-buildpack command", func() {
	var (
		requirementsFactory *testreq.FakeReqFactory
		repo                *testapi.FakeBuildpackRepository
		bitsRepo            *testapi.FakeBuildpackBitsRepository
		ui                  *testterm.FakeUI
		cmd                 CreateBuildpack
	)

	BeforeEach(func() {
		requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: true}
		repo = &testapi.FakeBuildpackRepository{}
		bitsRepo = &testapi.FakeBuildpackBitsRepository{}
		ui = &testterm.FakeUI{}
		cmd = NewCreateBuildpack(ui, repo, bitsRepo)
	})

	It("fails requirements when the user is not logged in", func() {
		requirementsFactory.LoginSuccess = false
		context := testcmd.NewContext("create-buildpack", []string{"my-buildpack", "my-dir", "0"})
		testcmd.RunCommand(cmd, context, requirementsFactory)

		Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
	})

	It("fails with usage when given fewer than three arguments", func() {
		context := testcmd.NewContext("create-buildpack", []string{})
		testcmd.RunCommand(cmd, context, requirementsFactory)

		Expect(ui.FailedWithUsage).To(BeTrue())
	})

	It("creates and uploads buildpacks", func() {
		context := testcmd.NewContext("create-buildpack", []string{"my-buildpack", "my.war", "5"})
		testcmd.RunCommand(cmd, context, requirementsFactory)

		Expect(repo.CreateBuildpack.Enabled).To(BeNil())
		Expect(ui.FailedWithUsage).To(BeFalse())

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Creating buildpack", "my-buildpack"},
			[]string{"OK"},
			[]string{"Uploading buildpack", "my-buildpack"},
			[]string{"OK"},
		))
		Expect(ui.Outputs).ToNot(ContainSubstrings([]string{"FAILED"}))
	})

	It("warns the user when the buildpack already exists", func() {
		repo.CreateBuildpackExists = true
		context := testcmd.NewContext("create-buildpack", []string{"my-buildpack", "my.war", "5"})
		testcmd.RunCommand(cmd, context, requirementsFactory)

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Creating buildpack", "my-buildpack"},
			[]string{"OK"},
			[]string{"my-buildpack", "already exists"},
			[]string{"TIP", "use", cf.Name(), "update-buildpack"},
		))
		Expect(ui.Outputs).ToNot(ContainSubstrings([]string{"FAILED"}))
	})

	It("enables the buildpack when given the --enabled flag", func() {
		context := testcmd.NewContext("create-buildpack", []string{"--enable", "my-buildpack", "my.war", "5"})
		testcmd.RunCommand(cmd, context, requirementsFactory)

		Expect(*repo.CreateBuildpack.Enabled).To(Equal(true))
	})

	It("disables the buildpack when given the --disable flag", func() {
		context := testcmd.NewContext("create-buildpack", []string{"--disable", "my-buildpack", "my.war", "5"})
		testcmd.RunCommand(cmd, context, requirementsFactory)
		Expect(*repo.CreateBuildpack.Enabled).To(Equal(false))
	})

	It("alerts the user when uploading the buildpack bits fails", func() {
		bitsRepo.UploadBuildpackErr = true
		context := testcmd.NewContext("create-buildpack", []string{"my-buildpack", "bogus/path", "5"})
		testcmd.RunCommand(cmd, context, requirementsFactory)

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Creating buildpack", "my-buildpack"},
			[]string{"OK"},
			[]string{"Uploading buildpack"},
			[]string{"FAILED"},
		))
	})
})
