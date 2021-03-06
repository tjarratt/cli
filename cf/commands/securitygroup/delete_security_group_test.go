package securitygroup_test

import (
	fakeSecurityGroup "github.com/cloudfoundry/cli/cf/api/security_groups/fakes"
	"github.com/cloudfoundry/cli/cf/configuration"
	"github.com/cloudfoundry/cli/cf/errors"
	"github.com/cloudfoundry/cli/cf/models"
	testcmd "github.com/cloudfoundry/cli/testhelpers/commands"
	testconfig "github.com/cloudfoundry/cli/testhelpers/configuration"
	testreq "github.com/cloudfoundry/cli/testhelpers/requirements"
	testterm "github.com/cloudfoundry/cli/testhelpers/terminal"

	. "github.com/cloudfoundry/cli/cf/commands/securitygroup"
	. "github.com/cloudfoundry/cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("delete-security-group command", func() {
	var (
		ui                  *testterm.FakeUI
		securityGroupRepo   *fakeSecurityGroup.FakeSecurityGroup
		requirementsFactory *testreq.FakeReqFactory
		configRepo          configuration.ReadWriter
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		requirementsFactory = &testreq.FakeReqFactory{}
		securityGroupRepo = &fakeSecurityGroup.FakeSecurityGroup{}
		configRepo = testconfig.NewRepositoryWithDefaults()
	})

	runCommand := func(args ...string) {
		cmd := NewDeleteSecurityGroup(ui, configRepo, securityGroupRepo)
		testcmd.RunCommand(cmd, args, requirementsFactory)
	}

	Describe("requirements", func() {
		It("should fail if not logged in", func() {
			runCommand("my-group")
			Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
		})

		It("should fail with usage when not provided a single argument", func() {
			requirementsFactory.LoginSuccess = true
			runCommand("whoops", "I can't believe", "I accidentally", "the whole thing")
			Expect(ui.FailedWithUsage).To(BeTrue())
		})
	})

	Context("when logged in", func() {
		BeforeEach(func() {
			requirementsFactory.LoginSuccess = true
		})

		Context("when the group with the given name exists", func() {
			BeforeEach(func() {
				securityGroupRepo.ReadReturns.SecurityGroup = models.SecurityGroup{
					SecurityGroupFields: models.SecurityGroupFields{
						Name: "my-group",
						Guid: "group-guid",
					},
				}
			})

			It("should delete the security group", func() {
				runCommand("my-group")
				Expect(securityGroupRepo.ReadCalledWith.Name).To(Equal("my-group"))
				Expect(securityGroupRepo.DeleteCalledWith.Guid).To(Equal("group-guid"))
			})

			It("tells the user what it's about to do", func() {
				runCommand("my-group")
				Expect(ui.Outputs).To(ContainSubstrings(
					[]string{"Deleting", "security group", "my-group", "my-user"},
					[]string{"OK"},
				))
			})
		})

		Context("when finding the group returns an error", func() {
			BeforeEach(func() {
				securityGroupRepo.ReadReturns.Error = errors.New("pbbbbbbbbbbt")
			})

			It("fails and tells the user", func() {
				runCommand("whoops")

				Expect(ui.Outputs).To(ContainSubstrings([]string{"FAILED"}))
			})
		})

		Context("when a group with that name does not exist", func() {
			BeforeEach(func() {
				securityGroupRepo.ReadReturns.Error = errors.NewModelNotFoundError("Security group", "uh uh uh -- you didn't sahy the magick word")
			})

			It("fails and tells the user", func() {
				runCommand("whoop")

				Expect(ui.WarnOutputs).To(ContainSubstrings([]string{"whoop", "does not exist"}))
			})
		})

		It("fails and warns the user if deleting fails", func() {
			securityGroupRepo.DeleteReturns.Error = errors.New("raspberry")
			runCommand("whoops")

			Expect(ui.Outputs).To(ContainSubstrings([]string{"FAILED"}))
		})
	})
})
