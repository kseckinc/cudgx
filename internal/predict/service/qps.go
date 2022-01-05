package service

import (
	"github.com/galaxy-future/cudgx/internal/predict/consts"
	"github.com/galaxy-future/cudgx/internal/predict/query"
)

//RedundancySeries 系统冗余度
type RedundancySeries struct {
	//ServiceName  服务名称
	ServiceName string `json:"service_name"`
	//MetricName指标名称
	MetricName string `json:"metric_name"`
	//Clusters 集群负载
	Clusters []*ClusterRedundancySeries `json:"clusters"`
}

//ClusterRedundancySeries 服务所属集群的冗余度
type ClusterRedundancySeries struct {
	//ClusterName 集群名称
	ClusterName string `json:"cluster"`
	//时间戳字段
	Timestamps []int64 `json:"timestamps"`
	//值字段
	Values []float64 `json:"values"`
}

//QueryRedundancyByQPS 基于QPS查询系统冗余度
func QueryRedundancyByQPS(serviceName, clusterName string, benchmark float64, begin, end int64) (*RedundancySeries, error) {
	samples, err := query.AverageQPS(serviceName, clusterName, begin, end)
	if err != nil {
		return nil, err
	}
	clusters := samples2ClusterSeries(samples)
	for _, cluster := range clusters {
		for i := range cluster.Values {
			cluster.Values[i] = benchmark / cluster.Values[i]
		}
	}
	series := &RedundancySeries{
		ServiceName: serviceName,
		MetricName:  consts.QPSMetricsName,
	}

	for _, cluster := range clusters {
		series.Clusters = append(series.Clusters, cluster)
	}
	return series, nil
}

//QueryServiceTotalQPS 基于QPS查询系统冗余度
func QueryServiceTotalQPS(serviceName, clusterName string, begin, end int64) (*RedundancySeries, error) {
	samples, err := query.TotalQPS(serviceName, clusterName, begin, end)
	if err != nil {
		return nil, err
	}

	clusters := samples2ClusterSeries(samples)

	series := &RedundancySeries{
		ServiceName: serviceName,
		MetricName:  consts.QPSMetricsName,
	}

	for _, cluster := range clusters {
		series.Clusters = append(series.Clusters, cluster)
	}
	return series, nil
}

func samples2ClusterSeries(samples []query.ClusterSample) []*ClusterRedundancySeries {
	clustersNameMap := make(map[string]*ClusterRedundancySeries)
	for _, sample := range samples {
		cluster := clustersNameMap[sample.ClusterName]
		if cluster == nil {
			cluster = &ClusterRedundancySeries{
				ClusterName: sample.ClusterName,
			}
			clustersNameMap[sample.ClusterName] = cluster
		}
		cluster.Values = append(cluster.Values, sample.Value)
		cluster.Timestamps = append(cluster.Timestamps, sample.Timestamp)
	}

	var clusters []*ClusterRedundancySeries
	for _, cluster := range clustersNameMap {
		clusters = append(clusters, cluster)
	}
	return clusters
}
