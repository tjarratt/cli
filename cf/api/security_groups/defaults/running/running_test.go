package running_test

import (
	"net/http"
	"net/http/httptest"
	"time"

	testapi "github.com/cloudfoundry/cli/cf/api/fakes"
	"github.com/cloudfoundry/cli/cf/configuration"
	"github.com/cloudfoundry/cli/cf/models"
	"github.com/cloudfoundry/cli/cf/net"
	testconfig "github.com/cloudfoundry/cli/testhelpers/configuration"
	testnet "github.com/cloudfoundry/cli/testhelpers/net"

	. "github.com/cloudfoundry/cli/cf/api/security_groups/defaults/running"
	. "github.com/cloudfoundry/cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RunningSecurityGroupsRepo", func() {
	var (
		testServer  *httptest.Server
		testHandler *testnet.TestHandler
		configRepo  configuration.ReadWriter
		repo        RunningSecurityGroupsRepo
	)

	BeforeEach(func() {
		configRepo = testconfig.NewRepositoryWithDefaults()
		gateway := net.NewCloudControllerGateway((configRepo), time.Now)
		repo = NewRunningSecurityGroupsRepo(configRepo, gateway)
	})

	AfterEach(func() {
		testServer.Close()
	})

	setupTestServer := func(reqs ...testnet.TestRequest) {
		testServer, testHandler = testnet.NewServer(reqs)
		configRepo.SetApiEndpoint(testServer.URL)
	}

	Describe(".AddToRunningSet", func() {
		It("makes a correct request", func() {
			setupTestServer(
				testapi.NewCloudControllerTestRequest(testnet.TestRequest{
					Method: "PUT",
					Path:   "/v2/config/running_security_groups/a-real-guid",
					Response: testnet.TestResponse{
						Status: http.StatusCreated,
						Body:   addRunningResponse,
					},
				}),
			)

			err := repo.AddToRunningSet("a-real-guid")

			Expect(err).ToNot(HaveOccurred())
			Expect(testHandler).To(HaveAllRequestsCalled())
		})
	})

	Describe("RemoveFromRunningSet", func() {
		It("makes a correct request", func() {
			testServer, testHandler = testnet.NewServer([]testnet.TestRequest{
				testapi.NewCloudControllerTestRequest(testnet.TestRequest{
					Method: "DELETE",
					Path:   "/v2/config/running_security_groups/my-guid",
					Response: testnet.TestResponse{
						Status: http.StatusNoContent,
					},
				}),
			})

			configRepo.SetApiEndpoint(testServer.URL)
			err := repo.RemoveFromRunningSet("my-guid")

			Expect(err).ToNot(HaveOccurred())
			Expect(testHandler).To(HaveAllRequestsCalled())
		})
	})

	Describe(".List", func() {
		It("returns a list of security groups that are the defaults for running", func() {
			setupTestServer(
				testapi.NewCloudControllerTestRequest(testnet.TestRequest{
					Method: "GET",
					Path:   "/v2/config/running_security_groups",
					Response: testnet.TestResponse{
						Status: http.StatusOK,
						Body:   firstRunningListItem,
					},
				}),
				testapi.NewCloudControllerTestRequest(testnet.TestRequest{
					Method: "GET",
					Path:   "/v2/config/running_security_groups",
					Response: testnet.TestResponse{
						Status: http.StatusOK,
						Body:   secondRunningListItem,
					},
				}),
			)

			defaults, err := repo.List()

			Expect(err).ToNot(HaveOccurred())
			Expect(testHandler).To(HaveAllRequestsCalled())
			Expect(defaults).To(Equal([]models.SecurityGroupFields{
				{
					Name: "name-71",
					Guid: "cd186158-b356-474d-9861-724f34f48502",
					Rules: []map[string]string{{
						"protocol": "udp",
					}},
				},
				{
					Name: "name-72",
					Guid: "d3374b62-7eac-4823-afbd-460d2bf44c67",
					Rules: []map[string]string{{
						"destination": "198.41.191.47/1",
					}},
				},
			}))
		})
	})
})

var addRunningResponse string = `{
  "metadata": {
    "guid": "897341eb-ef31-406f-b57b-414f51583a3a",
    "url": "/v2/config/running_security_groups/897341eb-ef31-406f-b57b-414f51583a3a",
    "created_at": "2014-06-23T21:43:30+00:00",
    "updated_at": "2014-06-23T21:43:30+00:00"
  },
  "entity": {
    "name": "name-904",
    "rules": [
      {
        "protocol": "udp",
        "ports": "8080",
        "destination": "198.41.191.47/1"
      }
    ]
  }
}`

var firstRunningListItem string = `{
  "next_url": "/v2/config/running_security_groups?page=2",
  "resources": [
    {
      "metadata": {
        "guid": "cd186158-b356-474d-9861-724f34f48502",
        "url": "/v2/security_groups/cd186158-b356-474d-9861-724f34f48502",
        "created_at": "2014-06-23T22:55:30+00:00",
        "updated_at": null
      },
      "entity": {
        "name": "name-71",
        "rules": [
          {
            "protocol": "udp"
          }
        ],
        "spaces_url": "/v2/security_groups/d3374b62-7eac-4823-afbd-460d2bf44c67/spaces"
      }
    }
  ]
}`

var secondRunningListItem string = `{
  "next_url": null,
  "resources": [
    {
      "metadata": {
        "guid": "d3374b62-7eac-4823-afbd-460d2bf44c67",
        "url": "/v2/config/running_security_groups/d3374b62-7eac-4823-afbd-460d2bf44c67",
        "created_at": "2014-06-23T22:55:30+00:00",
        "updated_at": null
      },
      "entity": {
        "name": "name-72",
        "rules": [
          {
            "destination": "198.41.191.47/1"
          }
        ],
        "spaces_url": "/v2/security_groups/d3374b62-7eac-4823-afbd-460d2bf44c67/spaces"
      }
    }
  ]
}`
