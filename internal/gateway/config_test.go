package gateway_test

import (
	"fmt"
	"github.com/galaxy-future/cudgx/common/kafka"
	"github.com/galaxy-future/cudgx/common/types"
	"github.com/galaxy-future/cudgx/internal/gateway"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"reflect"
	"time"
)

var _ = ginkgo.Describe("Config", func() {
	type configEntry struct {
		fileName string
		config   *gateway.Config
	}

	var (
		conf configEntry
	)

	ginkgo.BeforeEach(func() {
		conf.fileName = "tests/01-configure.json"
		conf.config = &gateway.Config{
			MonitoringRoute: &gateway.MessageRouteConfig{
				Entries: []*gateway.StorageEntryConfig{
					{
						ServicePrefix: "monitoring-gf",
						Brokers:       []string{"monitoring-gf:9092"},
						Topic:         "monitoring_gf_test",
					},
				},
				Default: &gateway.StorageEntryConfig{
					ServicePrefix: "",
					Brokers:       []string{"monitoring-default:9092"},
					Topic:         "monitoring_default_test",
				},
			},
			StreamingRoute: &gateway.MessageRouteConfig{
				Entries: []*gateway.StorageEntryConfig{
					{
						ServicePrefix: "streaming-gf",
						Brokers:       []string{"streaming-gf:9092"},
						Topic:         "streaming_gf_test",
					},
				},
				Default: &gateway.StorageEntryConfig{
					ServicePrefix: "",
					Brokers:       []string{"streaming-default:9092"},
					Topic:         "streaming_default_test",
				},
			},
			Producer: &kafka.ProducerConfig{
				MaxMessageBytes:  123456,
				RequiredAcks:     "WaitForLocal",
				Timeout:          types.Duration{Duration: 10 * time.Second},
				Compression:      "none",
				CompressionLevel: -1000,
				Return: struct {
					Successes bool
					Errors    bool
				}{false, true},
				Flush: struct {
					Bytes       int
					Messages    int
					Frequency   types.Duration
					MaxMessages int
				}{2089, 100, types.Duration{Duration: 4 * time.Second}, 2000},
				Retry: struct {
					Max     int
					Backoff types.Duration
				}{100, types.Duration{Duration: 600 * time.Millisecond}},
			},
		}
	})

	ginkgo.Describe("Load configure from file", func() {

		ginkgo.Context("Regular config file", func() {
			ginkgo.It("Load from file", func() {
				fmt.Println(conf.fileName)
				gateway, err := gateway.NewFromConfigFile(conf.fileName)
				//current , _:= json.Marshal(gateway.GetConfig())
				//expect , _ := json.Marshal(conf.config)
				//fmt.Println(string(current))
				//fmt.Println(string(expect))

				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(reflect.DeepEqual(gateway.GetConfig(), conf.config)).To(gomega.Equal(true))
			})
		})
	})
})
