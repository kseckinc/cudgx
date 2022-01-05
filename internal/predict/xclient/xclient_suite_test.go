package xclient_test

import (
	"github.com/galaxy-future/cudgx/internal/predict/xclient"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func TestXclient(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Xclient Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	xclient.InitializeBridgxClient("http://bridgx-api.internal.galaxy-future.org")
	xclient.InitializeSchedulxClient("http://10.16.23.96:9090")
})
