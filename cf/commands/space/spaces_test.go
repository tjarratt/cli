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

package space_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tjarratt/cli/cf/api"
	. "github.com/tjarratt/cli/cf/commands/space"
	"github.com/tjarratt/cli/cf/configuration"
	"github.com/tjarratt/cli/cf/models"
	testapi "github.com/tjarratt/cli/testhelpers/api"
	testcmd "github.com/tjarratt/cli/testhelpers/commands"
	testconfig "github.com/tjarratt/cli/testhelpers/configuration"
	testreq "github.com/tjarratt/cli/testhelpers/requirements"
	testterm "github.com/tjarratt/cli/testhelpers/terminal"

	. "github.com/tjarratt/cli/testhelpers/matchers"
)

func callSpaces(args []string, requirementsFactory *testreq.FakeReqFactory, config configuration.Reader, spaceRepo api.SpaceRepository) (ui *testterm.FakeUI) {
	ui = new(testterm.FakeUI)
	ctxt := testcmd.NewContext("spaces", args)

	cmd := NewListSpaces(ui, config, spaceRepo)
	testcmd.RunCommand(cmd, ctxt, requirementsFactory)
	return
}

var _ = Describe("Testing with ginkgo", func() {

	It("TestSpacesRequirements", func() {
		spaceRepo := &testapi.FakeSpaceRepository{}
		config := testconfig.NewRepository()

		requirementsFactory := &testreq.FakeReqFactory{LoginSuccess: true, TargetedOrgSuccess: true}
		callSpaces([]string{}, requirementsFactory, config, spaceRepo)
		Expect(testcmd.CommandDidPassRequirements).To(BeTrue())

		requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: true, TargetedOrgSuccess: false}
		callSpaces([]string{}, requirementsFactory, config, spaceRepo)
		Expect(testcmd.CommandDidPassRequirements).To(BeFalse())

		requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: false, TargetedOrgSuccess: true}
		callSpaces([]string{}, requirementsFactory, config, spaceRepo)
		Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
	})

	It("TestListingSpaces", func() {
		space := models.Space{}
		space.Name = "space1"
		space2 := models.Space{}
		space2.Name = "space2"
		space3 := models.Space{}
		space3.Name = "space3"
		spaceRepo := &testapi.FakeSpaceRepository{
			Spaces: []models.Space{space, space2, space3},
		}

		config := testconfig.NewRepositoryWithDefaults()
		requirementsFactory := &testreq.FakeReqFactory{LoginSuccess: true, TargetedOrgSuccess: true}

		ui := callSpaces([]string{}, requirementsFactory, config, spaceRepo)

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Getting spaces in org", "my-org", "my-user"},
			[]string{"space1"},
			[]string{"space2"},
			[]string{"space3"},
		))
	})

	It("TestListingSpacesWhenNoSpaces", func() {
		spaceRepo := &testapi.FakeSpaceRepository{
			Spaces: []models.Space{},
		}

		configRepo := testconfig.NewRepositoryWithDefaults()
		requirementsFactory := &testreq.FakeReqFactory{LoginSuccess: true, TargetedOrgSuccess: true}

		ui := callSpaces([]string{}, requirementsFactory, configRepo, spaceRepo)

		Expect(ui.Outputs).To(ContainSubstrings(
			[]string{"Getting spaces in org", "my-org", "my-user"},
			[]string{"No spaces found"},
		))
	})
})
