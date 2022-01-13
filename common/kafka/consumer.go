package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/galaxy-future/cudgx/common/logger"
	"go.uber.org/zap"
)

func NewConsumers(ch chan interface{}, brokers []string, topic, group string, kafkaConfig *ConsumerConfig) (*ConsumerClient, error) {
	if len(brokers) == 0 {
		return nil, fmt.Errorf("no Kafka bootstrap brokers defined")
	}
	if len(topic) == 0 {
		return nil, fmt.Errorf("no Kafka topic defined")
	}
	if len(group) == 0 {
		return nil, fmt.Errorf("no Kafka consumer group defined")
	}

	config := sarama.NewConfig()
	err := applyKafkaConfigure(config, kafkaConfig)
	if err != nil {
		return nil, err
	}

	client, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
	}

	return &ConsumerClient{
		client:    client,
		messageCh: ch,
		topic:     topic,
		readyCh:   make(chan struct{}),
	}, nil

}

func (client *ConsumerClient) Start(ctx context.Context) {
	for {
		if err := client.client.Consume(ctx, []string{client.topic}, client); err != nil {
			logger.GetLogger().Error("Error from consumer ", zap.String("error", err.Error()))
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			break
		}
		client.readyCh = make(chan struct{})
	}
}

func (client *ConsumerClient) Stop() {
	client.client.Close()
}

type ConsumerClient struct {
	client    sarama.ConsumerGroup
	topic     string
	readyCh   chan struct{}
	messageCh chan interface{}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (client *ConsumerClient) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(client.readyCh)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (client *ConsumerClient) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (client *ConsumerClient) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		client.messageCh <- message.Value
		session.MarkMessage(message, "")
	}

	return nil
}

func applyKafkaConfigure(saramaConf *sarama.Config, kafkaConfig *ConsumerConfig) error {

	version, err := sarama.ParseKafkaVersion(kafkaConfig.KafkaVersion)
	if err != nil {
		return err
	} else {
		saramaConf.Version = version
	}
	if kafkaConfig.Group.Session.Timeout.Duration != 0 {
		saramaConf.Consumer.Group.Session.Timeout = kafkaConfig.Group.Session.Timeout.Duration
	}

	if kafkaConfig.Group.Heartbeat.Interval.Duration != 0 {
		saramaConf.Consumer.Group.Heartbeat.Interval = kafkaConfig.Group.Heartbeat.Interval.Duration
	}

	switch kafkaConfig.Group.Rebalance.Strategy {
	case "range":
		saramaConf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	case "roundrobin":
		saramaConf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	default:
		saramaConf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	}
	if kafkaConfig.Group.Rebalance.Timeout.Duration != 0 {
		saramaConf.Consumer.Group.Rebalance.Timeout = kafkaConfig.Group.Rebalance.Timeout.Duration
	}
	if kafkaConfig.Group.Rebalance.Retry.Max == 0 {
		saramaConf.Consumer.Group.Rebalance.Retry.Max = 4
	} else {
		saramaConf.Consumer.Group.Rebalance.Retry.Max = kafkaConfig.Group.Rebalance.Retry.Max
	}
	if kafkaConfig.Group.Rebalance.Retry.Backoff.Duration != 0 {
		saramaConf.Consumer.Group.Rebalance.Retry.Backoff = kafkaConfig.Group.Rebalance.Retry.Backoff.Duration
	}

	if kafkaConfig.Retry.Backoff.Duration != 0 {
		saramaConf.Consumer.Retry.Backoff = kafkaConfig.Retry.Backoff.Duration
	}

	if kafkaConfig.Fetch.Min == 0 {
		saramaConf.Consumer.Fetch.Min = 1
	} else {
		saramaConf.Consumer.Fetch.Min = kafkaConfig.Fetch.Min
	}

	if kafkaConfig.Fetch.Default == 0 {
		saramaConf.Consumer.Fetch.Default = 1000000
	} else {
		saramaConf.Consumer.Fetch.Default = kafkaConfig.Fetch.Default
	}

	if kafkaConfig.Fetch.Max == 0 {
		saramaConf.Consumer.Fetch.Max = 0
	} else {
		saramaConf.Consumer.Fetch.Max = kafkaConfig.Fetch.Max
	}

	if kafkaConfig.MaxWaitTime.Duration != 0 {
		saramaConf.Consumer.MaxWaitTime = kafkaConfig.MaxWaitTime.Duration
	}

	if kafkaConfig.MaxProcessingTime.Duration != 0 {
		saramaConf.Consumer.MaxProcessingTime = kafkaConfig.MaxProcessingTime.Duration
	}

	saramaConf.Consumer.Return.Errors = kafkaConfig.Return.Errors

	switch kafkaConfig.Offsets.Initial {
	case "oldest":
		saramaConf.Consumer.Offsets.Initial = sarama.OffsetOldest
	case "newest":
		saramaConf.Consumer.Offsets.Initial = sarama.OffsetNewest
	default:
		saramaConf.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	if kafkaConfig.Offsets.CommitInterval.Duration != 0 {
		saramaConf.Consumer.Offsets.CommitInterval = kafkaConfig.Offsets.CommitInterval.Duration
	}

	if kafkaConfig.Offsets.Retention.Duration != 0 {
		saramaConf.Consumer.Offsets.Retention = kafkaConfig.Offsets.Retention.Duration
	}

	if kafkaConfig.Offsets.Retry.Max == 0 {
		saramaConf.Consumer.Offsets.Retry.Max = 3
	} else {
		saramaConf.Consumer.Offsets.Retry.Max = kafkaConfig.Offsets.Retry.Max
	}

	return nil
}
