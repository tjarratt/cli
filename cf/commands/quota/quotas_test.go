package quota_test

import (
	. "github.com/tjarratt/cli/cf/commands/quota"
	"github.com/tjarratt/cli/cf/models"
	. "github.com/tjarratt/cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/tjarratt/cli/cf/errors"
	testapi "github.com/tjarratt/cli/testhelpers/api"
	testcmd "github.com/tjarratt/cli/testhelpers/commands"
	testconfig "github.com/tjarratt/cli/testhelpers/configuration"
	testreq "github.com/tjarratt/cli/testhelpers/requirements"
	testterm "github.com/tjarratt/cli/testhelpers/terminal"
)

var _ = Describe("quotas command", func() {
	var (
		ui                  *testterm.FakeUI
		quotaRepo           *testapi.FakeQuotaRepository
		requirementsFactory *testreq.FakeReqFactory
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		quotaRepo = &testapi.FakeQuotaRepository{}
		requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: true}
	})

	runCommand := func() bool {
		cmd := NewListQuotas(ui, testconfig.NewRepositoryWithDefaults(), quotaRepo)
		return testcmd.RunCommand(cmd, testcmd.NewContext("create-quota", []string{}), requirementsFactory)
	}

	Describe("requirements", func() {
		It("requires the user to be logged in", func() {
			requirementsFactory.LoginSuccess = false
			Expect(runCommand()).ToNot(HavePassedRequirements())
		})
	})

	Context("when quotas exist", func() {
		BeforeEach(func() {
			quotaRepo.FindAllReturns.Quotas = []models.QuotaFields{
				models.QuotaFields{
					Name:                    "quota-name",
					MemoryLimit:             1024,
					RoutesLimit:             111,
					ServicesLimit:           222,
					NonBasicServicesAllowed: true,
				},
				models.QuotaFields{
					Name:                    "quota-non-basic-not-allowed",
					MemoryLimit:             434,
					RoutesLimit:             1,
					ServicesLimit:           2,
					NonBasicServicesAllowed: false,
				},
			}
		})

		It("lists quotas", func() {
			Expect(Expect(runCommand()).To(HavePassedRequirements())).To(HavePassedRequirements())
			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Getting quotas as", "my-user"},
				[]string{"OK"},
				[]string{"name", "memory limit", "routes", "service instances", "paid service plans"},
				[]string{"quota-name", "1G", "111", "222", "allowed"},
				[]string{"quota-non-basic-not-allowed", "434M", "1", "2", "disallowed"},
			))
		})
	})

	Context("when an error occurs fetching quotas", func() {
		BeforeEach(func() {
			quotaRepo.FindAllReturns.Error = errors.New("I haz a borken!")
		})

		It("prints an error", func() {
			Expect(runCommand()).To(HavePassedRequirements())
			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Getting quotas as", "my-user"},
				[]string{"FAILED"},
			))
		})
	})

})
