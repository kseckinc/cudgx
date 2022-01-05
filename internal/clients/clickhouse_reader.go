package clients

import (
	"database/sql"
	"fmt"
	"github.com/galaxy-future/cudgx/common/clickhouse"
)

var (
	ClickhouseRdCli *ClickhouseReader
)

func InitClickhouseRdCli(option *clickhouse.Config) (err error) {
	ClickhouseRdCli, err = NewClickhouse(option)
	return
}

type ClickhouseReader struct {
	Client   *sql.DB
	Database string
	Table    string
}

func NewClickhouse(option *clickhouse.Config) (*ClickhouseReader, error) {
	client, err := createConnection(option)
	if err != nil {
		return nil, err
	}
	return &ClickhouseReader{
		Client:   client,
		Database: option.Database,
		Table:    option.Table,
	}, nil
}

func createConnection(option *clickhouse.Config) (*sql.DB, error) {
	if len(option.Hosts) == 0 {
		return nil, fmt.Errorf("can not create connection , hosts can not be empty")
	}
	var connections []*sql.DB
	host := option.Hosts[0]
	dsn := fmt.Sprintf("%s://%s:%s@%s/%s?write_timeout=%s",
		option.Schema,
		option.User,
		option.Password,
		host,
		option.Database,
		option.WriteTimeout)

	connection, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return nil, err
	}
	err = connection.Ping()
	if err != nil {
		return nil, err
	}
	connections = append(connections, connection)

	return connection, nil
}
