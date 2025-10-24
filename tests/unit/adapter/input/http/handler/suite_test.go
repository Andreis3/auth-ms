//go:build unit

package handler_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func Test_CreateAuthUserHandlerSuite(t *testing.T) {
	suiteConfig, reporterConfig := GinkgoConfiguration()

	suiteConfig.SkipStrings = []string{"SKIPPED", "PENDING", "NEVER-RUN", "SKIP"}
	reporterConfig.FullTrace = true
	reporterConfig.Verbose = false

	RegisterFailHandler(Fail)
	RunSpecs(t, "CreateAuthUserHandler Suite Tests Context", suiteConfig, reporterConfig)
}
