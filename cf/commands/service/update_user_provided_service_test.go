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

package service_test

import (
	"github.com/tjarratt/cli/cf/api"
	. "github.com/tjarratt/cli/cf/commands/service"
	"github.com/tjarratt/cli/cf/models"
	testapi "github.com/tjarratt/cli/testhelpers/api"
	testcmd "github.com/tjarratt/cli/testhelpers/commands"
	testconfig "github.com/tjarratt/cli/testhelpers/configuration"
	testreq "github.com/tjarratt/cli/testhelpers/requirements"
	testterm "github.com/tjarratt/cli/testhelpers/terminal"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/tjarratt/cli/testhelpers/matchers"
)

func callUpdateUserProvidedService(args []string, requirementsFactory *testreq.FakeReqFactory, userProvidedServiceInstanceRepo api.UserProvidedServiceInstanceRepository) (fakeUI *testterm.FakeUI) {
	fakeUI = &testterm.FakeUI{}
	ctxt := testcmd.NewContext("update-user-provided-service", args)

	config := testconfig.NewRepositoryWithDefaults()

	cmd := NewUpdateUserProvidedService(fakeUI, config, userProvidedServiceInstanceRepo)
	testcmd.RunCommand(cmd, ctxt, requirementsFactory)
	return
}

var _ = Describe("Testing with ginkgo", func() {
	It("TestUpdateUserProvidedServiceFailsWithUsage", func() {
		requirementsFactory := &testreq.FakeReqFactory{}
		userProvidedServiceInstanceRepo := &testapi.FakeUserProvidedServiceInstanceRepo{}

		ui := callUpdateUserProvidedService([]string{}, requirementsFactory, userProvidedServiceInstanceRepo)
		Expect(ui.FailedWithUsage).To(BeTrue())

		ui = callUpdateUserProvidedService([]string{"foo"}, requirementsFactory, userProvidedServiceInstanceRepo)
		Expect(ui.FailedWithUsage).To(BeFalse())
	})

	It("TestUpdateUserProvidedServiceRequirements", func() {
		args := []string{"service-name"}
		requirementsFactory := &testreq.FakeReqFactory{}
		userProvidedServiceInstanceRepo := &testapi.FakeUserProvidedServiceInstanceRepo{}

		requirementsFactory.LoginSuccess = false
		callUpdateUserProvidedService(args, requirementsFactory, userProvidedServiceInstanceRepo)
		Expect(testcmd.CommandDidPassRequirements).To(BeFalse())

		requirementsFactory.LoginSuccess = true
		callUpdateUserProvidedService(args, requirementsFactory, userProvidedServiceInstanceRepo)
		Expect(testcmd.CommandDidPassRequirements).To(BeTrue())

		Expect(requirementsFactory.ServiceInstanceName).To(Equal("service-name"))
	})

	It("TestUpdateUserProvidedServiceWhenNoFlagsArePresent", func() {

		args := []string{"service-name"}
		serviceInstance := models.ServiceInstance{}
		serviceInstance.Name = "found-service-name"
		requirementsFactory := &testreq.FakeReqFactory{
			LoginSuccess:    true,
			ServiceInstance: serviceInstance,
		}
		repo := &testapi.FakeUserProvidedServiceInstanceRepo{}
		ui := callUpdateUserProvidedService(args, requirementsFactory, repo)

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Updating user provided service", "found-service-name", "my-org", "my-space", "my-user"},
			[]string{"OK"},
			[]string{"No changes"},
		))
	})

	It("TestUpdateUserProvidedServiceWithJson", func() {

		args := []string{"-p", `{"foo":"bar"}`, "-l", "syslog://example.com", "service-name"}
		serviceInstance := models.ServiceInstance{}
		serviceInstance.Name = "found-service-name"
		requirementsFactory := &testreq.FakeReqFactory{
			LoginSuccess:    true,
			ServiceInstance: serviceInstance,
		}
		repo := &testapi.FakeUserProvidedServiceInstanceRepo{}
		ui := callUpdateUserProvidedService(args, requirementsFactory, repo)

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Updating user provided service", "found-service-name", "my-org", "my-space", "my-user"},
			[]string{"OK"},
			[]string{"TIP"},
		))
		Expect(repo.UpdateServiceInstance.Name).To(Equal(serviceInstance.Name))
		Expect(repo.UpdateServiceInstance.Params).To(Equal(map[string]string{"foo": "bar"}))
		Expect(repo.UpdateServiceInstance.SysLogDrainUrl).To(Equal("syslog://example.com"))
	})

	It("TestUpdateUserProvidedServiceWithoutJson", func() {

		args := []string{"-l", "syslog://example.com", "service-name"}
		serviceInstance := models.ServiceInstance{}
		serviceInstance.Name = "found-service-name"
		requirementsFactory := &testreq.FakeReqFactory{
			LoginSuccess:    true,
			ServiceInstance: serviceInstance,
		}
		repo := &testapi.FakeUserProvidedServiceInstanceRepo{}
		ui := callUpdateUserProvidedService(args, requirementsFactory, repo)

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Updating user provided service"},
			[]string{"OK"},
		))
	})

	It("TestUpdateUserProvidedServiceWithInvalidJson", func() {

		args := []string{"-p", `{"foo":"ba`, "service-name"}
		serviceInstance := models.ServiceInstance{}
		serviceInstance.Name = "found-service-name"
		requirementsFactory := &testreq.FakeReqFactory{
			LoginSuccess:    true,
			ServiceInstance: serviceInstance,
		}
		userProvidedServiceInstanceRepo := &testapi.FakeUserProvidedServiceInstanceRepo{}

		ui := callUpdateUserProvidedService(args, requirementsFactory, userProvidedServiceInstanceRepo)

		Expect(userProvidedServiceInstanceRepo.UpdateServiceInstance).NotTo(Equal(serviceInstance))

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"FAILED"},
			[]string{"JSON is invalid"},
		))
	})

	It("TestUpdateUserProvidedServiceWithAServiceInstanceThatIsNotUserProvided", func() {

		args := []string{"-p", `{"foo":"bar"}`, "service-name"}
		plan := models.ServicePlanFields{}
		plan.Guid = "my-plan-guid"
		serviceInstance := models.ServiceInstance{}
		serviceInstance.Name = "found-service-name"
		serviceInstance.ServicePlan = plan

		requirementsFactory := &testreq.FakeReqFactory{
			LoginSuccess:    true,
			ServiceInstance: serviceInstance,
		}
		userProvidedServiceInstanceRepo := &testapi.FakeUserProvidedServiceInstanceRepo{}

		ui := callUpdateUserProvidedService(args, requirementsFactory, userProvidedServiceInstanceRepo)

		Expect(userProvidedServiceInstanceRepo.UpdateServiceInstance).NotTo(Equal(serviceInstance))

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"FAILED"},
			[]string{"Service Instance is not user provided"},
		))
	})
})
