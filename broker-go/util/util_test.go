package util_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cholick/cf-examples/broker-go/util"
)

var _ = Describe("Util", func() {
	Describe("GetPort", func() {
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
})
