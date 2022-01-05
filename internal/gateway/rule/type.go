package rule

import "time"

//type Metric struct {
//	Id                  int64
//	ServiceName         string
//	MetricName          string
//	MinInstanceCount    int
//	PredictSamplesCount int
//	Redundancy          float64
//	Updated             int64
//	Status              string
//	RuleId              int64
//}

type Rule struct {
	Id          int64
	ServiceName string
	MetricName  string
	Aggregate   string
	Filters     string
	Groups      string
	Benchmark   float64
	Ts          time.Time
}

//type MetricRule struct {
//	ServiceName         string
//	MetricName          string
//	MinInstanceCount    int
//	PredictSamplesCount int
//	Redundancy          float64
//	Updated             int64
//	Status              string
//
//	Aggregate string
//	Filters   string
//	Groups    string
//	Benchmark float64
//}

type Filter struct {
	Key    string
	Value  string
	Action string
}

type Aggregate struct {
	Operation string
	Param     string
}

type MysqlOption struct {
	Dsn            string `json:"dsn"`
	RefreshSeconds int    `json:"refresh_seconds"`
}
