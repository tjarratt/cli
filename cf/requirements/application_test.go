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
	"github.com/tjarratt/cli/cf/errors"
	"github.com/tjarratt/cli/cf/models"
	. "github.com/tjarratt/cli/cf/requirements"
	testapi "github.com/tjarratt/cli/testhelpers/api"
	testassert "github.com/tjarratt/cli/testhelpers/assert"
	testterm "github.com/tjarratt/cli/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testing with ginkgo", func() {

	It("TestApplicationReqExecute", func() {
		app := models.Application{}
		app.Name = "my-app"
		app.Guid = "my-app-guid"
		appRepo := &testapi.FakeApplicationRepository{}
		appRepo.ReadReturns.App = app
		ui := new(testterm.FakeUI)

		appReq := NewApplicationRequirement("foo", ui, appRepo)
		success := appReq.Execute()

		Expect(success).To(BeTrue())
		Expect(appRepo.ReadArgs.Name).To(Equal("foo"))
		Expect(appReq.GetApplication()).To(Equal(app))
	})

	It("TestApplicationReqExecuteWhenApplicationNotFound", func() {
		appRepo := &testapi.FakeApplicationRepository{}
		appRepo.ReadReturns.Error = errors.NewModelNotFoundError("app", "foo")
		ui := new(testterm.FakeUI)

		appReq := NewApplicationRequirement("foo", ui, appRepo)

		testassert.AssertPanic(testterm.FailedWasCalled, func() {
			appReq.Execute()
		})
	})
})
