package main

import (
	"github.com/galaxy-future/cudgx/common/logger"
	"github.com/galaxy-future/cudgx/internal/consumer"
	"context"
	"flag"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	configFile = flag.String("gf.cudgx.consumer.config", "conf/consumer.json", "consumer configure file")
)

func main() {
	flag.Parse()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger.GetLogger())

	file, err := os.Open(*configFile)
	if err != nil {
		panic("can not open configure file : " + *configFile)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic("read configure file failed : " + err.Error())
	}

	config, err := consumer.LoadConfig(data)
	if err != nil {
		panic(err)
	}

	csm, err := consumer.NewConsumer(config)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		csm.Start(ctx)
	}()

	//wait
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	//wait a signal
	<-c

	cancel()

	wg.Wait()
}
