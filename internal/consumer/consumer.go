package consumer

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/galaxy-future/cudgx/common/clickhouse"
	"github.com/galaxy-future/cudgx/common/kafka"
	"github.com/galaxy-future/cudgx/common/logger"
	"github.com/galaxy-future/cudgx/common/mod"
	"github.com/golang/protobuf/proto"
	clickhouseGo "github.com/mailru/go-clickhouse"
	"go.uber.org/zap"
)

type Consumer struct {
	kafkaClient      *kafka.ConsumerClient
	clickhouseWriter *clickhouse.AsyncWriter
	messageChan      chan interface{}
	config           *Config
}

func NewConsumer(config *Config) (*Consumer, error) {
	messagesCh := make(chan interface{}, 100000)
	consumer := &Consumer{
		config:      config,
		messageChan: messagesCh,
	}

	kafkaClient, err := kafka.NewConsumers(messagesCh, config.Kafka.Brokers, config.Kafka.Topic, config.Kafka.Group, config.Kafka.Consumer)
	if err != nil {
		return nil, err
	}
	consumer.kafkaClient = kafkaClient
	writer, err := clickhouse.NewWriter(config.Clickhouse, config.WriteConfig, messagesCh, consumer.commit)
	if err != nil {
		return nil, err
	}
	consumer.clickhouseWriter = writer

	return consumer, nil
}

func (consumer *Consumer) Start(ctx context.Context) {
	var wgKafka sync.WaitGroup
	wgKafka.Add(1)
	go func() {
		defer wgKafka.Done()
		consumer.kafkaClient.Start(ctx)
		logger.GetLogger().Info("kafka process exists")
	}()

	var wgWriter sync.WaitGroup
	wgWriter.Add(1)
	go func() {
		defer wgWriter.Done()
		consumer.clickhouseWriter.Start()
	}()
	<-ctx.Done()

	wgKafka.Wait()
	consumer.kafkaClient.Stop()

	close(consumer.messageChan)

	wgWriter.Wait()

}

func (consumer *Consumer) commit(connection *sql.DB, messages []interface{}) error {
	tx, err := connection.Begin()
	if err != nil {
		logger.GetLogger().Error("begin tx failed ", zap.Error(err))
		return err
	}
	stmt, err := tx.Prepare(fmt.Sprintf(`
		INSERT INTO %s.%s (
			metricName,
			serviceName,
			clusterName,
			serviceRegion,
			serviceAz,
			serviceHost,
			labelKeys,
			labelValues,
			timestamp,
		    value                            
		) VALUES (
			?, ?, ?, ?, ?,? ,? ,?,toDateTime(?), ? 
		)`, consumer.config.Clickhouse.Database, consumer.config.Clickhouse.Table))

	if err != nil {
		logger.GetLogger().Error("prepare stmt failed", zap.Error(err))
		return err
	}

	for _, data := range messages {
		binaryData, ok := data.([]byte)
		if !ok {
			logger.GetLogger().Error("message format error, can not convert interface data to []byte")
			continue
		}
		var batch mod.MetricBatch
		err := proto.Unmarshal(binaryData, &batch)
		if err != nil {
			logger.GetLogger().Error("unmarshal MetricBatch failed", zap.Error(err))
			continue
		}

		for _, metric := range batch.Messages {
			var keys, values []string
			for key, value := range metric.Labels {
				keys = append(keys, key)
				values = append(values, value)
			}

			if _, err := stmt.Exec(
				metric.MetricName,
				metric.ServiceName,
				metric.ClusterName,
				metric.ServiceRegion,
				metric.ServiceAz,
				metric.ServiceHost,
				clickhouseGo.Array(keys),
				clickhouseGo.Array(values),
				metric.Timestamp/1000,
				metric.Value,
			); err != nil {
				logger.GetLogger().Error("prepare stmt failed", zap.Error(err))
				return err
			}
		}
	}
	if err := tx.Commit(); err != nil {
		logger.GetLogger().Error("prepare stmt failed", zap.Error(err))
		return err
	}
	return nil
}
