package main

import (
	"flag"
	"fmt"
	"github.com/galaxy-future/cudgx/common/logger"
	metricGo "github.com/galaxy-future/metrics-go"
	"github.com/galaxy-future/metrics-go/aggregate"
	"github.com/galaxy-future/metrics-go/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math"
	"net"
	"time"
)

var (
	serverBind     = flag.String("gf.cudgx.sample.pi.bind", "0.0.0.0:8090", "server bind address default(0.0.0.0:8090)")
	goRoutineCount = flag.Int("gf.cudgx.sample.pi.count", 5000, "go routine count to calc pi")
)

var (
	latencyMin types.Metrics
	latencyMax types.Metrics
	latency    types.Metrics
	qps        types.Metrics
)

func main() {
	flag.Parse()
	r := gin.New()
	if gin.IsDebugging() {
		r.Use(gin.Logger())
	}
	r.Use(gin.Recovery())

	r.GET("/pi", HandlerCalcPi)
	r.GET("/", func(c *gin.Context) {
		c.String(200, "success")
	})

	l, err := net.Listen("tcp", *serverBind)
	if err != nil {
		logger.GetLogger().Error("server run failed ", zap.Error(err))
		panic("server listen failed ")
	}

	initMetrics()

	err = r.RunListener(l)
	if err != nil {
		logger.GetLogger().Error("server run failed ", zap.Error(err))
		panic("server start failed ")
	}

}

func HandlerCalcPi(c *gin.Context) {
	begin := time.Now()
	c.String(200, fmt.Sprintf("%v", pi(*goRoutineCount)))

	cost := time.Now().Sub(begin).Milliseconds()
	latencyMin.With().Value(float64(cost))
	latencyMax.With().Value(float64(cost))
	latency.With().Value(float64(cost))
	qps.With().Value(1)
}

// pi launches n goroutines to compute an
// approximation of pi.
func pi(n int) float64 {
	ch := make(chan float64)
	for k := 0; k < n; k++ {
		go term(ch, float64(k))
	}
	f := 0.0
	for k := 0; k < n; k++ {
		f += <-ch
	}
	return f
}

func term(ch chan float64, k float64) {
	ch <- 4 * math.Pow(-1, k) / (2*k + 1)
}

func initMetrics() {
	latencyMin = metricGo.NewMonitoringMetric("latencyMin", []string{}, aggregate.NewMinBuilder())
	latencyMax = metricGo.NewMonitoringMetric("latencyMax", []string{}, aggregate.NewMaxBuilder())
	latency = metricGo.NewStreamingMetric("latency", []string{})
	qps = metricGo.NewMonitoringMetric("qps", []string{}, aggregate.NewCountBuilder())
}
