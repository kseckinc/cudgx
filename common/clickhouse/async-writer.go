package clickhouse

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/galaxy-future/cudgx/common/logger"
	_ "github.com/mailru/go-clickhouse"
	"go.uber.org/zap"
)

type CommitFunc func(connection *sql.DB, messages []interface{}) error

//AsyncWriter 是clickhouse的写入器
type AsyncWriter struct {
	//connections 所有clickhouse的写入连接，clickhouse可以直接写入local表
	connections []*sql.DB
	//currentConnection 当前写入节点的index
	currentConnection int
	//MessageChan 写入缓存
	messagesCh <-chan interface{}
	//flush 写入配置
	config *WriterConfig

	//commitFunc 在每次消息累计到batch或者到指定时间时回调，数据写入到clickhouse
	commit CommitFunc
}

//NewWriter 新建一个clickhouse AsyncWriter
func NewWriter(config *Config, flush *WriterConfig, messagesCh <-chan interface{}, commit CommitFunc) (*AsyncWriter, error) {
	connections, err := createConnection(config)
	if err != nil {
		return nil, err
	}

	return &AsyncWriter{
		connections:       connections,
		currentConnection: 0,
		config:            flush,
		messagesCh:        messagesCh,
		commit:            commit,
	}, nil
}

func (writer *AsyncWriter) Start() {
	var wg sync.WaitGroup
	for i := 0; i < writer.config.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s := sender{
				writer:   writer,
				messages: make([]interface{}, 0, writer.config.BatchSize),
			}
			s.start()
		}()
	}

	wg.Wait()
}

type sender struct {
	writer            *AsyncWriter
	messages          []interface{}
	currentConnection int
}

func (s *sender) start() {
	ticker := time.NewTicker(s.writer.config.FlushDuration.Duration)
	for {
		select {
		case <-ticker.C:
			s.flushWithRetry(s.messages)
			s.messages = s.messages[:0]
		case item, ok := <-s.writer.messagesCh:
			if ok {
				s.messages = append(s.messages, item)
				if len(s.messages) == s.writer.config.BatchSize {
					s.flushWithRetry(s.messages)
					s.messages = s.messages[:0]
				}
			} else {
				s.flushWithRetry(s.messages)
				return
			}
		}
	}
}
func (s *sender) flushWithRetry(messages []interface{}) {
	err := s.flush(messages)
	if err == nil {
		return
	}
	logger.GetLogger().Error("failed when send metrics", zap.Error(err))
	time.Sleep(s.writer.config.Backoff.Duration)
	for i := 0; i < s.writer.config.RetryCount; i++ {
		err := s.flush(messages)
		if err == nil {
			return
		}
		logger.GetLogger().Error("failed when send metrics", zap.Error(err))
		time.Sleep(s.writer.config.Backoff.Duration)
	}
}

func (s *sender) flush(messages []interface{}) error {
	if s.currentConnection >= len(s.writer.connections) {
		s.currentConnection = 0
	}
	connection := s.writer.connections[s.currentConnection]
	s.currentConnection++
	return s.writer.commit(connection, messages)
}

//createConnection 基于config创建多个连接
func createConnection(config *Config) ([]*sql.DB, error) {
	if len(config.Hosts) == 0 {
		return nil, fmt.Errorf("can not create connection , hosts can not be empty")
	}
	var connections []*sql.DB
	for _, host := range config.Hosts {
		dsn := fmt.Sprintf("%s://%s:%s@%s/%s?write_timeout=%s",
			config.Schema,
			config.User,
			config.Password,
			host,
			config.Database,
			config.WriteTimeout)

		connection, err := sql.Open("clickhouse", dsn)
		if err != nil {
			return nil, err
		}
		err = connection.Ping()
		if err != nil {
			return nil, err
		}
		connections = append(connections, connection)
	}

	return connections, nil
}
