package kafka

import (
	"github.com/Shopify/sarama"
)

type ProducerClient struct {
	client sarama.AsyncProducer
}

func (producer *ProducerClient) SendMessage(topicName, key string, data []byte) {
	producer.client.Input() <- &sarama.ProducerMessage{
		Topic: topicName,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(data),
	}
}

func NewProducer(brokers []string, config *ProducerConfig) (*ProducerClient, error) {
	saramaConfig := sarama.NewConfig()
	applyKafkaProducerConfig(config, saramaConfig)

	client, err := sarama.NewAsyncProducer(brokers, saramaConfig)
	if err != nil {
		return nil, err
	}

	return &ProducerClient{client: client}, nil
}

func applyKafkaProducerConfig(conf *ProducerConfig, saramaConfig *sarama.Config) {
	if conf.MaxMessageBytes == 0 {
		saramaConfig.Producer.MaxMessageBytes = 1000000
	} else {
		saramaConfig.Producer.MaxMessageBytes = conf.MaxMessageBytes
	}

	switch conf.RequiredAcks {
	case "WaitForLocal":
		saramaConfig.Producer.RequiredAcks = sarama.WaitForLocal
	case "WaitForAll":
		saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	case "NoResponse":
		saramaConfig.Producer.RequiredAcks = sarama.NoResponse
	default:
		saramaConfig.Producer.RequiredAcks = sarama.WaitForLocal
	}

	if conf.Timeout.Duration != 0 {
		saramaConfig.Producer.Timeout = conf.Timeout.Duration
	}

	switch conf.Compression {
	case "none":
		saramaConfig.Producer.Compression = sarama.CompressionNone
	case "gzip":
		saramaConfig.Producer.Compression = sarama.CompressionGZIP
	case "snappy":
		saramaConfig.Producer.Compression = sarama.CompressionSnappy
	case "lz4":
		saramaConfig.Producer.Compression = sarama.CompressionLZ4
	default:
		saramaConfig.Producer.Compression = sarama.CompressionNone
	}

	if conf.CompressionLevel == 0 {
		saramaConfig.Producer.CompressionLevel = -1000
	} else {
		saramaConfig.Producer.CompressionLevel = conf.CompressionLevel
	}

	saramaConfig.Producer.Return.Successes = conf.Return.Successes
	saramaConfig.Producer.Return.Errors = conf.Return.Errors

	if conf.Flush.Frequency.Duration != 0 {
		saramaConfig.Producer.Flush.Frequency = conf.Flush.Frequency.Duration
	}

	saramaConfig.Producer.Flush.Bytes = conf.Flush.Bytes
	saramaConfig.Producer.Flush.Messages = conf.Flush.Messages
	saramaConfig.Producer.Flush.MaxMessages = conf.Flush.MaxMessages

	saramaConfig.Producer.Retry.Max = conf.Retry.Max

	if conf.Retry.Backoff.Duration != 0 {
		saramaConfig.Producer.Retry.Backoff = conf.Retry.Backoff.Duration
	}

}
