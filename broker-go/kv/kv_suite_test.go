package kv_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestKv(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kv Suite")
}
