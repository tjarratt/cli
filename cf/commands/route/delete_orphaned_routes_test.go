package route_test

import (
	"github.com/tjarratt/cli/cf/models"
	testapi "github.com/tjarratt/cli/testhelpers/api"
	testcmd "github.com/tjarratt/cli/testhelpers/commands"
	testconfig "github.com/tjarratt/cli/testhelpers/configuration"
	testreq "github.com/tjarratt/cli/testhelpers/requirements"
	testterm "github.com/tjarratt/cli/testhelpers/terminal"

	. "github.com/tjarratt/cli/cf/commands/route"
	. "github.com/tjarratt/cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("delete-orphaned-routes command", func() {
	var routeRepo *testapi.FakeRouteRepository
	var reqFactory *testreq.FakeReqFactory

	BeforeEach(func() {
		routeRepo = &testapi.FakeRouteRepository{}
		reqFactory = &testreq.FakeReqFactory{}
	})

	It("fails requirements when not logged in", func() {
		callDeleteOrphanedRoutes("y", []string{}, reqFactory, routeRepo)
		Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
	})

	Context("when logged in successfully", func() {

		BeforeEach(func() {
			reqFactory.LoginSuccess = true
		})

		It("passes requirements when logged in", func() {
			callDeleteOrphanedRoutes("y", []string{}, reqFactory, routeRepo)
			Expect(testcmd.CommandDidPassRequirements).To(BeTrue())
		})

		It("passes when confirmation is provided", func() {
			var ui *testterm.FakeUI
			domain := models.DomainFields{Name: "example.com"}
			domain2 := models.DomainFields{Name: "cookieclicker.co"}

			app1 := models.ApplicationFields{Name: "dora"}

			route := models.Route{}
			route.Host = "hostname-1"
			route.Domain = domain
			route.Apps = []models.ApplicationFields{app1}

			route2 := models.Route{}
			route2.Guid = "route2-guid"
			route2.Host = "hostname-2"
			route2.Domain = domain2

			routeRepo.Routes = []models.Route{route, route2}

			ui = callDeleteOrphanedRoutes("y", []string{}, reqFactory, routeRepo)

			Expect(ui.Prompts).To(ContainSubstrings(
				[]string{"Really delete orphaned routes"},
			))

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Deleting route", "hostname-2.cookieclicker.co"},
				[]string{"OK"},
			))
			Expect(routeRepo.DeletedRouteGuids).To(ContainElement("route2-guid"))
		})

		It("passes when the force flag is used", func() {
			var ui *testterm.FakeUI
			domain := models.DomainFields{Name: "example.com"}
			domain2 := models.DomainFields{Name: "cookieclicker.co"}

			app1 := models.ApplicationFields{Name: "dora"}

			route := models.Route{}
			route.Host = "hostname-1"
			route.Domain = domain
			route.Apps = []models.ApplicationFields{app1}

			route2 := models.Route{}
			route2.Guid = "route2-guid"
			route2.Host = "hostname-2"
			route2.Domain = domain2

			routeRepo.Routes = []models.Route{route, route2}

			ui = callDeleteOrphanedRoutes("", []string{"-f"}, reqFactory, routeRepo)

			Expect(len(ui.Prompts)).To(Equal(0))

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Deleting route", "hostname-2.cookieclicker.co"},
				[]string{"OK"},
			))
			Expect(routeRepo.DeletedRouteGuids).To(ContainElement("route2-guid"))
		})
	})
})

func callDeleteOrphanedRoutes(confirmation string, args []string, reqFactory *testreq.FakeReqFactory, routeRepo *testapi.FakeRouteRepository) (ui *testterm.FakeUI) {
	ui = &testterm.FakeUI{Inputs: []string{confirmation}}
	ctxt := testcmd.NewContext("delete-orphaned-routes", args)
	configRepo := testconfig.NewRepositoryWithDefaults()

	cmd := NewDeleteOrphanedRoutes(ui, configRepo, routeRepo)

	testcmd.RunCommand(cmd, ctxt, reqFactory)
	return
}
