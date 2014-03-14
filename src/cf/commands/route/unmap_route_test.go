package route_test

import (
	. "cf/commands/route"
	"cf/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tjarratt/cli"
	testapi "testhelpers/api"
	testassert "testhelpers/assert"
	testcmd "testhelpers/commands"
	testconfig "testhelpers/configuration"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
)

var _ = Describe("Unmap Route Command", func() {
	It("TestUnmapRouteFailsWithUsage", func() {
		reqFactory := &testreq.FakeReqFactory{}
		routeRepo := &testapi.FakeRouteRepository{}

		ui := callUnmapRoute([]string{}, reqFactory, routeRepo)
		Expect(ui.FailedWithUsage).To(BeTrue())

		ui = callUnmapRoute([]string{"foo"}, reqFactory, routeRepo)
		Expect(ui.FailedWithUsage).To(BeTrue())

		ui = callUnmapRoute([]string{"foo", "bar"}, reqFactory, routeRepo)
		Expect(ui.FailedWithUsage).To(BeFalse())
	})

	It("TestUnmapRouteRequirements", func() {
		routeRepo := &testapi.FakeRouteRepository{}
		reqFactory := &testreq.FakeReqFactory{LoginSuccess: true}

		callUnmapRoute([]string{"-n", "my-host", "my-app", "my-domain.com"}, reqFactory, routeRepo)
		Expect(testcmd.CommandDidPassRequirements).To(BeTrue())
		Expect(reqFactory.ApplicationName).To(Equal("my-app"))
		Expect(reqFactory.DomainName).To(Equal("my-domain.com"))
	})

	It("TestUnmapRouteWhenUnbinding", func() {
		domain := models.DomainFields{
			Guid: "my-domain-guid",
			Name: "example.com",
		}
		route := models.Route{RouteSummary: models.RouteSummary{
			Domain: domain,
			RouteFields: models.RouteFields{
				Guid: "my-route-guid",
				Host: "foo",
			},
		}}
		app := models.Application{ApplicationFields: models.ApplicationFields{
			Guid: "my-app-guid",
			Name: "my-app",
		}}

		routeRepo := &testapi.FakeRouteRepository{FindByHostAndDomainRoute: route}
		reqFactory := &testreq.FakeReqFactory{LoginSuccess: true, Application: app, Domain: domain}

		ui := callUnmapRoute([]string{"-n", "my-host", "my-app", "my-domain.com"}, reqFactory, routeRepo)

		testassert.SliceContains(ui.Outputs, testassert.Lines{
			{"Removing route", "foo.example.com", "my-app", "my-org", "my-space", "my-user"},
			{"OK"},
		})

		Expect(routeRepo.UnboundRouteGuid).To(Equal("my-route-guid"))
		Expect(routeRepo.UnboundAppGuid).To(Equal("my-app-guid"))
	})
})

func callUnmapRoute(args []string, reqFactory *testreq.FakeReqFactory, routeRepo *testapi.FakeRouteRepository) (ui *testterm.FakeUI) {
	ui = new(testterm.FakeUI)
	var ctxt *cli.Context = testcmd.NewContext("unmap-route", args)

	configRepo := testconfig.NewRepositoryWithDefaults()
	cmd := NewUnmapRoute(ui, configRepo, routeRepo)
	testcmd.RunCommand(cmd, ctxt, reqFactory)
	return
}
