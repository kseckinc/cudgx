package kafka

import (
	"github.com/galaxy-future/cudgx/common/types"
)

//ConsumerConfig is from https://github.com/Shopify/sarama
type ConsumerConfig struct {
	KafkaVersion string
	// Group is the namespace for configuring consumer group.
	Group struct {
		Session struct {
			// The timeout used to detect consumer failures when using Kafka's group management facility.
			// The consumer sends periodic heartbeats to indicate its liveness to the broker.
			// If no heartbeats are received by the broker before the expiration of this session timeout,
			// then the broker will remove this consumer from the group and initiate a rebalance.
			// Note that the value must be in the allowable range as configured in the broker configuration
			// by `group.min.session.timeout.ms` and `group.max.session.timeout.ms` (default 10s)
			Timeout types.Duration
		}
		Heartbeat struct {
			// The expected time between heartbeats to the consumer coordinator when using Kafka's group
			// management facilities. Heartbeats are used to ensure that the consumer's session stays active and
			// to facilitate rebalancing when new consumers join or leave the group.
			// The value must be set lower than ConsumerClient.Group.Session.Timeout, but typically should be set no
			// higher than 1/3 of that value.
			// It can be adjusted even lower to control the expected time for normal rebalances (default 3s)
			Interval types.Duration
		}
		Rebalance struct {
			// Strategy for allocating topic partitions to members (default BalanceStrategyRange)
			Strategy string
			// The maximum allowed time for each worker to join the group once a rebalance has begun.
			// This is basically a limit on the amount of time needed for all tasks to flush any pending
			// data and commit offsets. If the timeout is exceeded, then the worker will be removed from
			// the group, which will cause offset commit failures (default 60s).
			Timeout types.Duration

			Retry struct {
				// When a new consumer joins a consumer group the set of consumers attempt to "rebalance"
				// the load to assign partitions to each consumer. If the set of consumers changes while
				// this assignment is taking place the rebalance will fail and retry. This setting controls
				// the maximum number of attempts before giving up (default 4).
				Max int
				// Backoff time between retries during rebalance (default 2s)
				Backoff types.Duration
			}
		}
		Member struct {
			// Custom metadata to include when joining the group. The user data for all joined members
			// can be retrieved by sending a DescribeGroupRequest to the broker that is the
			// coordinator for the group.
			UserData []byte
		}
	}

	Retry struct {
		// How long to wait after a failing to read from a partition before
		// trying again (default 2s).
		Backoff types.Duration
	}

	// Fetch is the namespace for controlling how many bytes are retrieved by any
	// given request.
	Fetch struct {
		// The minimum number of message bytes to fetch in a request - the broker
		// will wait until at least this many are available. The default is 1,
		// as 0 causes the consumer to spin when no messages are available.
		// Equivalent to the JVM's `fetch.min.bytes`.
		Min int32
		// The default number of message bytes to fetch from the broker in each
		// request (default 1MB). This should be larger than the majority of
		// your messages, or else the consumer will spend a lot of time
		// negotiating sizes and not actually consuming. Similar to the JVM's
		// `fetch.message.max.bytes`.
		Default int32
		// The maximum number of message bytes to fetch from the broker in a
		// single request. Messages larger than this will return
		// ErrMessageTooLarge and will not be consumable, so you must be sure
		// this is at least as large as your largest message. Defaults to 0
		// (no limit). Similar to the JVM's `fetch.message.max.bytes`. The
		// global `sarama.MaxResponseSize` still applies.
		Max int32
	}
	// The maximum amount of time the broker will wait for ConsumerClient.Fetch.Min
	// bytes to become available before it returns fewer than that anyways. The
	// default is 250ms, since 0 causes the consumer to spin when no events are
	// available. 100-500ms is a reasonable range for most cases. Kafka only
	// supports precision up to milliseconds; nanoseconds will be truncated.
	// Equivalent to the JVM's `fetch.wait.max.ms`.
	MaxWaitTime types.Duration

	// The maximum amount of time the consumer expects a message takes to
	// process for the user. If writing to the Messages channel takes longer
	// than this, that partition will stop fetching more messages until it
	// can proceed again.
	// Note that, since the Messages channel is buffered, the actual grace time is
	// (MaxProcessingTime * ChanneBufferSize). Defaults to 100ms.
	// If a message is not written to the Messages channel between two ticks
	// of the expiryTicker then a timeout is detected.
	// Using a ticker instead of a timer to detect timeouts should typically
	// result in many fewer calls to Timer functions which may result in a
	// significant performance improvement if many messages are being sent
	// and timeouts are infrequent.
	// The disadvantage of using a ticker instead of a timer is that
	// timeouts will be less accurate. That is, the effective timeout could
	// be between `MaxProcessingTime` and `2 * MaxProcessingTime`. For
	// example, if `MaxProcessingTime` is 100ms then a delay of 180ms
	// between two messages being sent may not be recognized as a timeout.
	MaxProcessingTime types.Duration

	// Return specifies what channels will be populated. If they are set to true,
	// you must read from them to prevent deadlock.
	Return struct {
		// If enabled, any errors that occurred while consuming are returned on
		// the Errors channel (default disabled).
		Errors bool
	}

	// Offsets specifies configuration for how and when to commit consumed
	// offsets. This currently requires the manual use of an OffsetManager
	// but will eventually be automated.
	Offsets struct {
		// How frequently to commit updated offsets. Defaults to 1s.
		CommitInterval types.Duration

		// The initial offset to use if no offset was previously committed.
		// Should be OffsetNewest or OffsetOldest. Defaults to OffsetNewest.
		Initial string

		// The retention duration for committed offsets. If zero, disabled
		// (in which case the `offsets.retention.minutes` option on the
		// broker will be used).  Kafka only supports precision up to
		// milliseconds; nanoseconds will be truncated. Requires Kafka
		// broker version 0.9.0 or later.
		// (default is 0: disabled).
		Retention types.Duration

		Retry struct {
			// The total number of times to retry failing commit
			// requests during OffsetManager shutdown (default 3).
			Max int
		}
	}
}

//ProducerConfig 配置Kafka Producer，从https://github.com/Shopify/sarama而来
type ProducerConfig struct {
	// The maximum permitted size of a message (defaults to 1000000). Should be
	// set equal to or smaller than the broker's `message.max.bytes`.
	MaxMessageBytes int
	// The level of acknowledgement reliability needed from the broker (defaults
	// to WaitForLocal). Equivalent to the `request.required.acks` setting of the
	// JVM producer.
	RequiredAcks string
	// The maximum duration the broker will wait the receipt of the number of
	// RequiredAcks (defaults to 10 seconds). This is only relevant when
	// RequiredAcks is set to WaitForAll or a number > 1. Only supports
	// millisecond resolution, nanoseconds will be truncated. Equivalent to
	// the JVM producer's `request.timeout.ms` setting.
	Timeout types.Duration
	// The type of compression to use on messages (defaults to no compression).
	// Similar to `compression.codec` setting of the JVM producer.
	Compression string
	// The level of compression to use on messages. The meaning depends
	// on the actual compression type used and defaults to default compression
	// level for the codec.
	CompressionLevel int

	// Return specifies what channels will be populated. If they are set to true,
	// you must read from the respective channels to prevent deadlock. If,
	// however, this config is used to create a `SyncProducer`, both must be set
	// to true and you shall not read from the channels since the producer does
	// this internally.
	Return struct {
		// If enabled, successfully delivered messages will be returned on the
		// Successes channel (default disabled).
		Successes bool

		// If enabled, messages that failed to deliver will be returned on the
		// Errors channel, including error (default enabled).
		Errors bool
	}

	// The following config options control how often messages are batched up and
	// sent to the broker. By default, messages are sent as fast as possible, and
	// all messages received while the current batch is in-flight are placed
	// into the subsequent batch.
	Flush struct {
		// The best-effort number of bytes needed to trigger a flush. Use the
		// global sarama.MaxRequestSize to set a hard upper limit.
		Bytes int
		// The best-effort number of messages needed to trigger a flush. Use
		// `MaxMessages` to set a hard upper limit.
		Messages int
		// The best-effort frequency of flushes. Equivalent to
		// `queue.buffering.max.ms` setting of JVM producer.
		Frequency types.Duration
		// The maximum number of messages the producer will send in a single
		// broker request. Defaults to 0 for unlimited. Similar to
		// `queue.buffering.max.messages` in the JVM producer.
		MaxMessages int
	}

	Retry struct {
		// The total number of times to retry sending a message (default 3).
		// Similar to the `message.send.max.retries` setting of the JVM producer.
		Max int
		// How long to wait for the cluster to settle between retries
		// (default 100ms). Similar to the `retry.backoff.ms` setting of the
		// JVM producer.
		Backoff types.Duration
	}
}
