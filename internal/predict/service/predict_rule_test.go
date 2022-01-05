package service

import (
	"github.com/galaxy-future/cudgx/internal/request"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("PredictRule", func() {

	ginkgo.Context("CreatePredictRule", func() {
		pr := &request.CreatePredictRuleRequest{
			Name:             "test-predict-rule",
			ServiceName:      "gf.sample.service",
			ClusterName:      "gf.cluster",
			MetricName:       "qps",
			BenchmarkQps:     100,
			MinRedundancy:    100,
			MaxRedundancy:    300,
			MinInstanceCount: 3,
			MaxInstanceCount: 10,
			ExecuteRatio:     30,
			Status:           "enable",
		}
		ginkgo.It("创建扩缩容规则", func() {
			err := CreatePredictRule(pr)
			gomega.Expect(err).To(gomega.BeNil())
		})
		ginkgo.It("查询扩缩容集群列表", func() {
			list, total, err := ListPredictRules("gf.sample.service", "gf.sample.service", 1, 20)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(total > 0).To(gomega.BeTrue())
			gomega.Expect(len(list) > 0).To(gomega.BeTrue())
		})
	})
})
