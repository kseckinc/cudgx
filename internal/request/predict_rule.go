package request

type CreatePredictRuleRequest struct {
	Name             string `json:"name" binding:"required"`
	ServiceName      string `json:"service_name" binding:"required"`
	ClusterName      string `json:"cluster_name" binding:"required"`
	MetricName       string `json:"metric_name" binding:"required"`
	BenchmarkQps     int    `json:"benchmark_qps" binding:"required"`
	MinRedundancy    int    `json:"min_redundancy" binding:"required"`
	MaxRedundancy    int    `json:"max_redundancy" binding:"required"`
	MinInstanceCount int    `json:"min_instance_count" binding:"required"`
	MaxInstanceCount int    `json:"max_instance_count" binding:"required"`
	ExecuteRatio     int    `json:"execute_ratio" binding:"required"`
	Status           string `json:"status" binding:"required"`
}

type UpdatePredictRuleRequest struct {
	Id               int64  `json:"id" binding:"required"`
	Name             string `json:"name" binding:"required"`
	ServiceName      string `json:"service_name" binding:"required"`
	ClusterName      string `json:"cluster_name" binding:"required"`
	MetricName       string `json:"metric_name" binding:"required"`
	BenchmarkQps     int    `json:"benchmark_qps" binding:"required"`
	MinRedundancy    int    `json:"min_redundancy" binding:"required"`
	MaxRedundancy    int    `json:"max_redundancy" binding:"required"`
	MinInstanceCount int    `json:"min_instance_count" binding:"required"`
	MaxInstanceCount int    `json:"max_instance_count" binding:"required"`
	ExecuteRatio     int    `json:"execute_ratio" binding:"required"`
	Status           string `json:"status" binding:"required"`
}

type BatchDeletePredictRuleRequest struct {
	Ids []int64 `json:"ids" binding:"min=1"`
}

type EnableOrDisablePredictRuleRequest struct {
	Id     int64  `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"`
}
