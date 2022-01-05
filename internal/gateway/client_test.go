package gateway_test

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	"github.com/galaxy-future/cudgx/internal/gateway"
)

var _ = ginkgo.Describe("Client", func() {
	var (
		g *gateway.Gateway
	)

	ginkgo.BeforeEach(func() {
		gateway, err := gateway.NewFromConfigFile("tests/02-client.json")
		if err != nil {
			ginkgo.Fail("failed create gateway")
		}
		g = gateway
	})

	ginkgo.Describe("Get configure entry ", func() {
		ginkgo.Context("Get service entry", func() {
			ginkgo.It("Get default monitoring entry", func() {
				entry := g.GetMonitoringStorageConfigEntry("aaaa", "bbbb")
				gomega.Expect(entry.ServicePrefix).To(gomega.Equal(""))
				gomega.Expect(entry.Topic).To(gomega.Equal("monitoring_metrics_test"))
			})

			ginkgo.It("Get gf monitoring entry", func() {
				entry := g.GetMonitoringStorageConfigEntry("gf.metrics.pi", "bbbb")
				gomega.Expect(entry.ServicePrefix).To(gomega.Equal("gf"))
				gomega.Expect(entry.Topic).To(gomega.Equal("monitoring_metrics_gf_test"))
			})

			ginkgo.It("Get default streaming entry", func() {
				entry := g.GetStreamingStorageConfigEntry("aaaa", "bbb")
				gomega.Expect(entry.ServicePrefix).To(gomega.Equal(""))
				gomega.Expect(entry.Topic).To(gomega.Equal("streaming_metrics_test"))
			})
			ginkgo.It("Get gf streaming entry", func() {
				entry := g.GetStreamingStorageConfigEntry("gf.metrics.pi", "bbb")
				gomega.Expect(entry.ServicePrefix).To(gomega.Equal("gf"))
				gomega.Expect(entry.Topic).To(gomega.Equal("streaming_metrics_gf_test"))
			})
		})
	})

	ginkgo.Describe("Get client", func() {
		ginkgo.Context("Get service client", func() {
			ginkgo.It("Get monitoring client", func() {
				clientAAA, err := g.GetMonitoringWriter("aaaa", "bbb")
				gomega.Expect(err).To(gomega.BeNil())
				clientBBB, err := g.GetMonitoringWriter("aaaa", "bbb")
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(clientAAA).To(gomega.Equal(clientBBB))

				clientGFPI, err := g.GetMonitoringWriter("gf.metrics.pi", "bbb")
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(clientAAA).NotTo(gomega.Equal(clientGFPI))

				clientGFConsumer, err := g.GetMonitoringWriter("gf.metrics.consumer", "bbb")
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(clientAAA).NotTo(gomega.Equal(clientGFConsumer))
				gomega.Expect(clientGFConsumer).To(gomega.Equal(clientGFPI))
			})

			ginkgo.It("Get streaming client", func() {
				clientAAA, err := g.GetStreamingWriter("aaaa", "bbb")
				gomega.Expect(err).To(gomega.BeNil())
				clientBBB, err := g.GetStreamingWriter("aaaa", "bbb")
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(clientAAA).To(gomega.Equal(clientBBB))

				clientGFPI, err := g.GetStreamingWriter("gf.metrics.pi", "bbb")
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(clientAAA).NotTo(gomega.Equal(clientGFPI))

				clientGFConsumer, err := g.GetStreamingWriter("gf.metrics.consumer", "bbb")
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(clientAAA).NotTo(gomega.Equal(clientGFConsumer))
				gomega.Expect(clientGFConsumer).To(gomega.Equal(clientGFPI))
			})

		})
	})

})
