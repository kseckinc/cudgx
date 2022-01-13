package service

import (
	"github.com/galaxy-future/cudgx/internal/predict/consts"
	"github.com/galaxy-future/cudgx/internal/predict/query"
)

//QueryInstances 基于QPS查询服务节点个数
func QueryInstances(serviceName, clusterName string, begin, end int64, trimmedSecond int64) (*RedundancySeries, error) {
	samples, err := query.InstanceCount(serviceName, clusterName, begin, end)
	if err != nil {
		return nil, err
	}

	clusters := samples2ClusterSeries(samples, trimmedSecond)

	series := &RedundancySeries{
		ServiceName: serviceName,
		MetricName:  consts.QPSMetricsName,
	}

	for _, cluster := range clusters {
		series.Clusters = append(series.Clusters, cluster)
	}
	return series, nil
}
