package thrifty_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestThrifty(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Thrifty Suite")
}
