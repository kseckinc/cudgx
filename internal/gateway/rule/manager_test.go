package rule_test

import (
	"github.com/galaxy-future/cudgx/common/mod"
	"github.com/galaxy-future/cudgx/internal/gateway/rule"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Manager", func() {
	validServiceName := "gf.cudgx.sample.pi"
	validMetricName := "latency"
	invalidServiceName := "unknown.service"
	invalidMetricName := "unknown.metric"

	ginkgo.Context("NewRuleManager", func() {
		var invalidOption *rule.MysqlOption

		ginkgo.BeforeEach(func() {
			invalidOption = &rule.MysqlOption{
				Dsn:            "root:wrongPasssword@tcp(localhost:3306)/cudgx_test?charset=utf8mb4&parseTime=True&loc=Local",
				RefreshSeconds: 5,
			}
		})

		ginkgo.It("Using an invalid dsn", func() {
			_, err := rule.NewRuleManager(invalidOption)
			gomega.Expect(err).NotTo(gomega.BeNil())
		})
	})

	ginkgo.Context("MatchRule", func() {
		ginkgo.It("When finding a exists servcie", func() {
			rule := manager.MatchRule(validServiceName, validMetricName)
			gomega.Expect(rule).NotTo(gomega.BeNil())
			gomega.Expect(rule.MetricName).To(gomega.Equal(validMetricName))
			gomega.Expect(rule.ServiceName).To(gomega.Equal(validServiceName))
		})
		ginkgo.It("When finding a missing service", func() {
			rule := manager.MatchRule(invalidServiceName, invalidMetricName)
			gomega.Expect(rule).To(gomega.BeNil())
		})
	})

	ginkgo.Context("WrapStreamingMessage", func() {
		ginkgo.It("When wrap a exist service metrics", func() {
			streamingBatch := mod.StreamingBatch{
				ServiceName: validServiceName,
				MetricName:  validMetricName,
				Messages: []*mod.StreamingMessage{
					{
						ServiceName: validServiceName,
						MetricName:  validMetricName,
					},
				},
			}
			wraped, err := manager.WrapStreamingMessage(&streamingBatch)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(len(wraped.Messages)).To(gomega.Equal(len(streamingBatch.GetMessages())))
			gomega.Expect(wraped.Messages[0].Rule.Aggregate.Operation).To(gomega.Equal("section_factor"))
			gomega.Expect(wraped.Messages[0].Rule.Aggregate.Param).To(gomega.Equal("{\"sections\":[10,30,40,50,100],\"factors\":[0.01,0.1,0.3,0.5,1,10]}"))
			gomega.Expect(len(wraped.Messages[0].Rule.Filters)).To(gomega.Equal(2))
			gomega.Expect(wraped.Messages[0].Rule.Benchmark).To(gomega.Equal(float64(500)))
		})
	})

})
