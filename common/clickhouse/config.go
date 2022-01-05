package clickhouse

import "github.com/galaxy-future/cudgx/common/types"

//Config consumer写入clickhouse配置
type Config struct {
	//Schema https/http
	Schema string `json:"schema"`
	//User clickhouse user
	User string `json:"user"`
	//Password clickhouse password
	Password string `json:"password"`
	//Database consumer消息写入的库
	Database string `json:"database"`
	//Table consumer 写入的表
	Table string `json:"table"`
	//Hosts clickhouse 节点
	Hosts []string `json:"hosts"`
	//WriteTimeout 写入超时
	WriteTimeout string `json:"write_timeout"`
	//ReadTimeout 查询超时
	ReadTimeout string `json:"read_timeout"`
}

type WriterConfig struct {
	FlushDuration types.Duration `json:"flush_duration"`
	RetryCount    int            `json:"retry_count"`
	Backoff       types.Duration `json:"backoff"`
	BatchSize     int            `json:"batch_size"`
	Concurrency   int            `json:"concurrency"`
}
