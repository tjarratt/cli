/*
                       WARNING WARNING WARNING

                Attention all potential contributors

   This testfile is not in the best state. We've been slowly transitioning
   from the built in "testing" package to using Ginkgo. As you can see, we've
   changed the format, but a lot of the setup, test body, descriptions, etc
   are either hardcoded, completely lacking, or misleading.

   For example:

   Describe("Testing with ginkgo"...)      // This is not a great description
   It("TestDoesSoemthing"...)              // This is a horrible description

   Describe("create-user command"...       // Describe the actual object under test
   It("creates a user when provided ..."   // this is more descriptive

   For good examples of writing Ginkgo tests for the cli, refer to

   src/github.com/tjarratt/cli/cf/commands/application/delete_app_test.go
   src/github.com/tjarratt/cli/cf/terminal/ui_test.go
   src/github.com/cloudfoundry/loggregator_consumer/consumer_test.go
*/

package requirements_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tjarratt/cli/cf/models"
	. "github.com/tjarratt/cli/cf/requirements"
	testassert "github.com/tjarratt/cli/testhelpers/assert"
	testconfig "github.com/tjarratt/cli/testhelpers/configuration"
	testterm "github.com/tjarratt/cli/testhelpers/terminal"

	. "github.com/tjarratt/cli/testhelpers/matchers"
)

var _ = Describe("Testing with ginkgo", func() {
	It("TestTargetedOrgRequirement", func() {
		ui := new(testterm.FakeUI)
		org := models.OrganizationFields{}
		org.Name = "my-org"
		org.Guid = "my-org-guid"
		config := testconfig.NewRepositoryWithDefaults()

		req := NewTargetedOrgRequirement(ui, config)
		success := req.Execute()
		Expect(success).To(BeTrue())

		config.SetOrganizationFields(models.OrganizationFields{})

		testassert.AssertPanic(testterm.FailedWasCalled, func() {
			NewTargetedOrgRequirement(ui, config).Execute()
		})

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"FAILED"},
			[]string{"No org targeted"},
		))
	})
})
