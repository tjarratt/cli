package requirements_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/errors"
	"github.com/tjarratt/cli/cf/models"
	. "github.com/tjarratt/cli/cf/requirements"
	testapi "github.com/tjarratt/cli/testhelpers/api"
	testassert "github.com/tjarratt/cli/testhelpers/assert"
	testconfig "github.com/tjarratt/cli/testhelpers/configuration"
	testterm "github.com/tjarratt/cli/testhelpers/terminal"
)

var _ = Describe("Testing with ginkgo", func() {
	var config configuration.ReadWriter
	var ui *testterm.FakeUI

	BeforeEach(func() {
		ui = new(testterm.FakeUI)
		config = testconfig.NewRepository()
		config.SetOrganizationFields(models.OrganizationFields{Guid: "the-org-guid"})
	})

	It("succeeds when the domain is found", func() {
		domain := models.DomainFields{Name: "example.com", Guid: "domain-guid"}
		domainRepo := &testapi.FakeDomainRepository{FindByNameInOrgDomain: domain}
		domainReq := NewDomainRequirement("example.com", ui, config, domainRepo)
		success := domainReq.Execute()

		Expect(success).To(BeTrue())
		Expect(domainRepo.FindByNameInOrgName).To(Equal("example.com"))
		Expect(domainRepo.FindByNameInOrgGuid).To(Equal("the-org-guid"))
		Expect(domainReq.GetDomain()).To(Equal(domain))
	})

	It("fails when the domain is not found", func() {
		domainRepo := &testapi.FakeDomainRepository{FindByNameInOrgApiResponse: errors.NewModelNotFoundError("Domain", "")}
		domainReq := NewDomainRequirement("example.com", ui, config, domainRepo)

		testassert.AssertPanic(testterm.FailedWasCalled, func() {
			domainReq.Execute()
		})
	})

	It("fails when an error occurs fetching the domain", func() {
		domainRepo := &testapi.FakeDomainRepository{FindByNameInOrgApiResponse: errors.NewWithError("", errors.New(""))}
		domainReq := NewDomainRequirement("example.com", ui, config, domainRepo)

		testassert.AssertPanic(testterm.FailedWasCalled, func() {
			domainReq.Execute()
		})
	})
})
