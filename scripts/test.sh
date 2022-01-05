#!/bin/bash

echo "starting mysql"
# start mysql
mysqld --initialize-insecure
mysqld --daemonize --pid-file=/run/mysqld/mysqld.pid --user=root &

echo "starting kafka"
# start kafka
cd /kafka/kafka_2.12-2.7.1
./bin/zookeeper-server-start.sh config/zookeeper.properties 2>1 1>zookeeper.log  &
sleep 5
./bin/kafka-server-start.sh  ./config/server.properties 2>1 1>kafka.log  &
sleep 5

bin/kafka-topics.sh --bootstrap-server=localhost:9092 --create --topic=test
sleep 5


echo "starting clickhouse"
chown -R root:root /var/lib/clickhouse
clickhouse-server --config=/etc/clickhouse-server/config.xml 2>1 1>clickhouse.log &


echo "begin test..."
cd go/src/github.com/galaxy-future/cudgx/
ginkgo ./...




