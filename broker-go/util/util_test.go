package util_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cholick/cf-examples/broker-go/util"
)

var _ = Describe("Util", func() {
	Describe("port", func() {
		It("returns default value", func() {
			port := util.GetPort()

			Expect(port).To(Equal("3000"))
		})

		It("pulls port from environment", func() {
			os.Setenv("PORT", "8765")
			port := util.GetPort()

			Expect(port).To(Equal("8765"))
		})
	})

	Describe("user", func() {
		It("returns default value", func() {
			port := util.GetUser()

			Expect(port).To(Equal("user"))
		})

		It("pulls from environment", func() {
			os.Setenv("SECURITY_USER_NAME", "admin")
			port := util.GetUser()

			Expect(port).To(Equal("admin"))
		})
	})

	Describe("user", func() {
		It("returns default value", func() {
			port := util.GetPassword()

			Expect(port).To(Equal("pass"))
		})

		It("pulls from environment", func() {
			os.Setenv("SECURITY_USER_PASSWORD", "secure")
			port := util.GetPassword()

			Expect(port).To(Equal("secure"))
		})
	})
})
