package models_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/tjarratt/cli/cf/models"
)

var _ = Describe("DomainFields", func() {
	var route Route

	BeforeEach(func() {
		route = Route{}

		domain := DomainFields{}
		domain.Name = "example.com"
		route.Domain = domain
	})

	It("uses the hostname as part of the URL", func() {
		route.Host = "foo"
		Expect(route.URL()).To(Equal("foo.example.com"))
	})

	It("omits the hostname when none is given", func() {
		Expect(route.URL()).To(Equal("example.com"))
	})
})
