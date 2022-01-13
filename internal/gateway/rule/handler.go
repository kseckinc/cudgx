package rule

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/galaxy-future/cudgx/common/mod"
	wrapmod "github.com/galaxy-future/cudgx/internal/gateway/mod"
)

func (manager *Manager) WrapStreamingMessage(streamingBatch *mod.StreamingBatch) (batch *wrapmod.StreamingRuleBatch, err error) {
	streamingMessages := streamingBatch.GetMessages()
	wrapStreamingBatch := wrapmod.StreamingRuleBatch{}
	wrapMetrics := make([]*wrapmod.StreamingRuleMessage, len(streamingMessages))
	for index, streamingMessage := range streamingMessages {
		rule := manager.MatchRule(streamingMessage.ServiceName, streamingMessage.MetricName)
		if rule == nil {
			continue
		}

		var filters []*wrapmod.Filter
		err = json.Unmarshal([]byte(rule.Filters), &filters)
		if err != nil {
			panic("wrap metric failed")
		}
		var aggregate *wrapmod.Aggregate
		fmt.Println(rule.Aggregate)
		err = json.Unmarshal([]byte(rule.Aggregate), &aggregate)
		if err != nil {
			panic("wrap metric failed")
		}

		metricRule := wrapmod.Rule{
			Benchmark: rule.Benchmark,
			Filters:   filters,
			Groups:    strings.Split(rule.Groups, ","),
			Aggregate: aggregate,
		}

		wrapMetric := wrapmod.StreamingRuleMessage{
			ServiceName:   streamingMessage.ServiceName,
			ServiceHost:   streamingMessage.ServiceHost,
			ServiceRegion: streamingMessage.ServiceRegion,
			ServiceAz:     streamingMessage.ServiceAz,
			MetricName:    streamingMessage.MetricName,
			ClusterName:   streamingMessage.ClusterName,
			Labels:        streamingMessage.Labels,
			Timestamp:     streamingMessage.Timestamp,
			Values:        streamingMessage.Values,
			Rule:          &metricRule,
		}
		wrapMetrics[index] = &wrapMetric
	}
	wrapStreamingBatch.MetricName = streamingBatch.MetricName
	wrapStreamingBatch.ServiceName = streamingBatch.ServiceName
	wrapStreamingBatch.Messages = wrapMetrics
	return &wrapStreamingBatch, err
}
