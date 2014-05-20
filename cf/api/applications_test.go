package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/tjarratt/cli/cf/errors"
	"github.com/tjarratt/cli/cf/models"
	"github.com/tjarratt/cli/cf/net"
	testapi "github.com/tjarratt/cli/testhelpers/api"
	testconfig "github.com/tjarratt/cli/testhelpers/configuration"
	testnet "github.com/tjarratt/cli/testhelpers/net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/tjarratt/cli/cf/api"
	. "github.com/tjarratt/cli/testhelpers/matchers"
)

var _ = Describe("ApplicationsRepository", func() {
	Describe("finding apps by name", func() {
		It("returns the app when it is found", func() {
			ts, handler, repo := createAppRepo([]testnet.TestRequest{findAppRequest})
			defer ts.Close()

			app, apiErr := repo.Read("My App")
			Expect(handler).To(testnet.HaveAllRequestsCalled())
			Expect(apiErr).NotTo(HaveOccurred())
			Expect(app.Name).To(Equal("My App"))
			Expect(app.Guid).To(Equal("app1-guid"))
			Expect(app.Memory).To(Equal(uint64(128)))
			Expect(app.DiskQuota).To(Equal(uint64(512)))
			Expect(app.InstanceCount).To(Equal(1))
			Expect(app.EnvironmentVars).To(Equal(map[string]string{"foo": "bar", "baz": "boom"}))
			Expect(app.Routes[0].Host).To(Equal("app1"))
			Expect(app.Routes[0].Domain.Name).To(Equal("cfapps.io"))
			Expect(app.Stack.Name).To(Equal("awesome-stacks-ahoy"))
		})

		It("returns a failure response when the app is not found", func() {
			request := testapi.NewCloudControllerTestRequest(findAppRequest)
			request.Response = testnet.TestResponse{Status: http.StatusOK, Body: `{"resources": []}`}

			ts, handler, repo := createAppRepo([]testnet.TestRequest{request})
			defer ts.Close()

			_, apiErr := repo.Read("My App")
			Expect(handler).To(testnet.HaveAllRequestsCalled())
			Expect(apiErr.(*errors.ModelNotFoundError)).NotTo(BeNil())
		})
	})

	Describe("creating applications", func() {
		It("makes the right request", func() {
			ts, handler, repo := createAppRepo([]testnet.TestRequest{createApplicationRequest})
			defer ts.Close()

			params := defaultAppParams()
			createdApp, apiErr := repo.Create(params)

			Expect(handler).To(testnet.HaveAllRequestsCalled())
			Expect(apiErr).NotTo(HaveOccurred())

			app := models.Application{}
			app.Name = "my-cool-app"
			app.Guid = "my-cool-app-guid"
			Expect(createdApp).To(Equal(app))
		})

		It("omits fields that are not set", func() {
			request := testapi.NewCloudControllerTestRequest(testnet.TestRequest{
				Method:   "POST",
				Path:     "/v2/apps",
				Matcher:  testnet.RequestBodyMatcher(`{"name":"my-cool-app","instances":3,"memory":2048,"disk_quota":512,"space_guid":"some-space-guid"}`),
				Response: testnet.TestResponse{Status: http.StatusCreated, Body: createApplicationResponse},
			})

			ts, handler, repo := createAppRepo([]testnet.TestRequest{request})
			defer ts.Close()

			params := defaultAppParams()
			params.BuildpackUrl = nil
			params.StackGuid = nil
			params.Command = nil

			_, apiErr := repo.Create(params)
			Expect(handler).To(testnet.HaveAllRequestsCalled())
			Expect(apiErr).NotTo(HaveOccurred())
		})
	})

	Describe("reading environment for an app", func() {
		Context("when the response can be parsed as json", func() {
			var (
				testServer   *httptest.Server
				userEnv      map[string]string
				vcapServices string
				err          error
				handler      *testnet.TestHandler
				repo         ApplicationRepository
			)

			AfterEach(func() {
				testServer.Close()
			})

			Context("when there are system provided env vars", func() {
				BeforeEach(func() {

					var appEnvRequest = testapi.NewCloudControllerTestRequest(testnet.TestRequest{
						Method: "GET",
						Path:   "/v2/apps/some-cool-app-guid/env",
						Response: testnet.TestResponse{
							Status: http.StatusOK,
							Body: `
{
   "system_env_json": {
      "VCAP_SERVICES": {
        "system_hash": {
          "system_key": "system_value"
        }
      }
   },
   "environment_json": {
      "key": "value"
   }
}
`,
						}})

					testServer, handler, repo = createAppRepo([]testnet.TestRequest{appEnvRequest})
					userEnv, vcapServices, err = repo.ReadEnv("some-cool-app-guid")
					Expect(err).ToNot(HaveOccurred())
					Expect(handler).To(testnet.HaveAllRequestsCalled())
				})

				It("returns the user environment and vcap services", func() {
					Expect(userEnv["key"]).To(Equal("value"))
					Expect(strings.Split(vcapServices, "\n")).To(ContainSubstrings(
						[]string{"system_hash", ":", "{"},
						[]string{"system_key", ":", "system_value"},
						[]string{"}"},
					))
				})
			})

			Context("when there are no environment variables", func() {
				BeforeEach(func() {
					var emptyEnvRequest = testapi.NewCloudControllerTestRequest(testnet.TestRequest{
						Method: "GET",
						Path:   "/v2/apps/some-cool-app-guid/env",
						Response: testnet.TestResponse{
							Status: http.StatusOK,
							Body:   `{"system_env_json": {"VCAP_SERVICES": {} }}`,
						}})

					testServer, handler, repo = createAppRepo([]testnet.TestRequest{emptyEnvRequest})
					userEnv, vcapServices, err = repo.ReadEnv("some-cool-app-guid")
					Expect(err).ToNot(HaveOccurred())
					Expect(handler).To(testnet.HaveAllRequestsCalled())
				})

				It("returns an empty string", func() {
					Expect(userEnv["key"]).To(BeEmpty())
					Expect(vcapServices).To(BeEmpty())
				})
			})
		})
	})

	Describe("updating applications", func() {
		It("makes the right request", func() {
			ts, handler, repo := createAppRepo([]testnet.TestRequest{updateApplicationRequest})
			defer ts.Close()

			app := models.Application{}
			app.Guid = "my-app-guid"
			app.Name = "my-cool-app"
			app.BuildpackUrl = "buildpack-url"
			app.Command = "some-command"
			app.Memory = 2048
			app.InstanceCount = 3
			app.Stack = &models.Stack{Guid: "some-stack-guid"}
			app.SpaceGuid = "some-space-guid"
			app.State = "started"
			app.DiskQuota = 512
			Expect(app.EnvironmentVars).To(BeNil())

			updatedApp, apiErr := repo.Update(app.Guid, app.ToParams())

			Expect(handler).To(testnet.HaveAllRequestsCalled())
			Expect(apiErr).NotTo(HaveOccurred())
			Expect(updatedApp.Name).To(Equal("my-cool-app"))
			Expect(updatedApp.Guid).To(Equal("my-cool-app-guid"))
		})

		It("sets environment variables", func() {
			request := testapi.NewCloudControllerTestRequest(testnet.TestRequest{
				Method:   "PUT",
				Path:     "/v2/apps/app1-guid",
				Matcher:  testnet.RequestBodyMatcher(`{"environment_json":{"DATABASE_URL":"mysql://example.com/my-db"}}`),
				Response: testnet.TestResponse{Status: http.StatusCreated},
			})

			ts, handler, repo := createAppRepo([]testnet.TestRequest{request})
			defer ts.Close()

			envParams := map[string]string{"DATABASE_URL": "mysql://example.com/my-db"}
			params := models.AppParams{EnvironmentVars: &envParams}

			_, apiErr := repo.Update("app1-guid", params)

			Expect(handler).To(testnet.HaveAllRequestsCalled())
			Expect(apiErr).NotTo(HaveOccurred())
		})

		It("can remove environment variables", func() {
			request := testapi.NewCloudControllerTestRequest(testnet.TestRequest{
				Method:   "PUT",
				Path:     "/v2/apps/app1-guid",
				Matcher:  testnet.RequestBodyMatcher(`{"environment_json":{}}`),
				Response: testnet.TestResponse{Status: http.StatusCreated},
			})

			ts, handler, repo := createAppRepo([]testnet.TestRequest{request})
			defer ts.Close()

			envParams := map[string]string{}
			params := models.AppParams{EnvironmentVars: &envParams}

			_, apiErr := repo.Update("app1-guid", params)

			Expect(handler).To(testnet.HaveAllRequestsCalled())
			Expect(apiErr).NotTo(HaveOccurred())
		})
	})

	It("deletes applications", func() {
		deleteApplicationRequest := testapi.NewCloudControllerTestRequest(testnet.TestRequest{
			Method:   "DELETE",
			Path:     "/v2/apps/my-cool-app-guid?recursive=true",
			Response: testnet.TestResponse{Status: http.StatusOK, Body: ""},
		})

		ts, handler, repo := createAppRepo([]testnet.TestRequest{deleteApplicationRequest})
		defer ts.Close()

		apiErr := repo.Delete("my-cool-app-guid")
		Expect(handler).To(testnet.HaveAllRequestsCalled())
		Expect(apiErr).NotTo(HaveOccurred())
	})
})

var singleAppResponse = testnet.TestResponse{
	Status: http.StatusOK,
	Body: `
{
  "resources": [
    {
      "metadata": {
        "guid": "app1-guid"
      },
      "entity": {
        "name": "My App",
        "environment_json": {
      		"foo": "bar",
      		"baz": "boom"
    	},
        "memory": 128,
        "instances": 1,
        "disk_quota": 512,
        "state": "STOPPED",
        "stack": {
			"metadata": {
				  "guid": "app1-route-guid"
				},
			"entity": {
			  "name": "awesome-stacks-ahoy"
			  }
        },
        "routes": [
      	  {
      	    "metadata": {
      	      "guid": "app1-route-guid"
      	    },
      	    "entity": {
      	      "host": "app1",
      	      "domain": {
      	      	"metadata": {
      	      	  "guid": "domain1-guid"
      	      	},
      	      	"entity": {
      	      	  "name": "cfapps.io"
      	      	}
      	      }
      	    }
      	  }
        ]
      }
    }
  ]
}`}

var findAppRequest = testapi.NewCloudControllerTestRequest(testnet.TestRequest{
	Method:   "GET",
	Path:     "/v2/spaces/my-space-guid/apps?q=name%3AMy+App&inline-relations-depth=1",
	Response: singleAppResponse,
})

var createApplicationResponse = `
{
    "metadata": {
        "guid": "my-cool-app-guid"
    },
    "entity": {
        "name": "my-cool-app"
    }
}`

var createApplicationRequest = testapi.NewCloudControllerTestRequest(testnet.TestRequest{
	Method: "POST",
	Path:   "/v2/apps",
	Matcher: testnet.RequestBodyMatcher(`{
		"name":"my-cool-app",
		"instances":3,
		"buildpack":"buildpack-url",
		"memory":2048,
		"disk_quota": 512,
		"space_guid":"some-space-guid",
		"stack_guid":"some-stack-guid",
		"command":"some-command"
	}`),
	Response: testnet.TestResponse{
		Status: http.StatusCreated,
		Body:   createApplicationResponse},
})

func defaultAppParams() models.AppParams {
	name := "my-cool-app"
	buildpackUrl := "buildpack-url"
	spaceGuid := "some-space-guid"
	stackGuid := "some-stack-guid"
	command := "some-command"
	memory := uint64(2048)
	diskQuota := uint64(512)
	instanceCount := 3

	return models.AppParams{
		Name:          &name,
		BuildpackUrl:  &buildpackUrl,
		SpaceGuid:     &spaceGuid,
		StackGuid:     &stackGuid,
		Command:       &command,
		Memory:        &memory,
		DiskQuota:     &diskQuota,
		InstanceCount: &instanceCount,
	}
}

var updateApplicationResponse = `
{
    "metadata": {
        "guid": "my-cool-app-guid"
    },
    "entity": {
        "name": "my-cool-app"
    }
}`

var updateApplicationRequest = testapi.NewCloudControllerTestRequest(testnet.TestRequest{
	Method:  "PUT",
	Path:    "/v2/apps/my-app-guid?inline-relations-depth=1",
	Matcher: testnet.RequestBodyMatcher(`{"name":"my-cool-app","instances":3,"buildpack":"buildpack-url","memory":2048,"disk_quota":512,"space_guid":"some-space-guid","state":"STARTED","stack_guid":"some-stack-guid","command":"some-command"}`),
	Response: testnet.TestResponse{
		Status: http.StatusOK,
		Body:   updateApplicationResponse},
})

func createAppRepo(requests []testnet.TestRequest) (ts *httptest.Server, handler *testnet.TestHandler, repo ApplicationRepository) {
	ts, handler = testnet.NewServer(requests)
	configRepo := testconfig.NewRepositoryWithDefaults()
	configRepo.SetApiEndpoint(ts.URL)
	gateway := net.NewCloudControllerGateway(configRepo, time.Now)
	repo = NewCloudControllerApplicationRepository(configRepo, gateway)
	return
}
