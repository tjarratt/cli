package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/tjarratt/cli/cf/api"
	"github.com/tjarratt/cli/cf/net"
	testapi "github.com/tjarratt/cli/testhelpers/api"
	testconfig "github.com/tjarratt/cli/testhelpers/configuration"
	testnet "github.com/tjarratt/cli/testhelpers/net"
	"net/http"
	"net/http/httptest"
	"time"
)

var _ = Describe("CloudControllerPasswordRepository", func() {
	It("updates your password", func() {
		req := testapi.NewCloudControllerTestRequest(testnet.TestRequest{
			Method:   "PUT",
			Path:     "/Users/my-user-guid/password",
			Matcher:  testnet.RequestBodyMatcher(`{"password":"new-password","oldPassword":"old-password"}`),
			Response: testnet.TestResponse{Status: http.StatusOK},
		})

		passwordUpdateServer, handler, repo := createPasswordRepo(req)
		defer passwordUpdateServer.Close()

		apiErr := repo.UpdatePassword("old-password", "new-password")
		Expect(handler).To(testnet.HaveAllRequestsCalled())
		Expect(apiErr).NotTo(HaveOccurred())
	})
})

func createPasswordRepo(req testnet.TestRequest) (passwordServer *httptest.Server, handler *testnet.TestHandler, repo PasswordRepository) {
	passwordServer, handler = testnet.NewServer([]testnet.TestRequest{req})

	configRepo := testconfig.NewRepositoryWithDefaults()
	configRepo.SetUaaEndpoint(passwordServer.URL)
	gateway := net.NewCloudControllerGateway(configRepo, time.Now)
	repo = NewCloudControllerPasswordRepository(configRepo, gateway)
	return
}
