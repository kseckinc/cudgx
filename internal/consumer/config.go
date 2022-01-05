package consumer

import (
	"encoding/json"
	"github.com/galaxy-future/cudgx/common/clickhouse"
	"github.com/galaxy-future/cudgx/common/kafka"
)

func LoadConfig(data []byte) (*Config, error) {
	var config Config
	err := json.Unmarshal(data, &config)
	return &config, err
}

//Config 是consumer的配置
type Config struct {
	//Kafka 配置
	Kafka *KafkaConfig `json:"kafka"`
	//Clickhouse 连接配置
	Clickhouse *clickhouse.Config `json:"clickhouse"`
	//WriteConfig 写入配置
	WriteConfig *clickhouse.WriterConfig `json:"write_config"`
}

//KafkaConfig 消费程序用到kafka的配置
type KafkaConfig struct {
	//Brokers kafka brokers
	Brokers []string `json:"brokers"`
	//Group kafka 消费Group
	Group string `json:"group"`
	//Topic 消费Topic
	Topic string `json:"topic"`
	//Consumer consumer配置
	Consumer *kafka.ConsumerConfig `json:"consumer"`
}
