create database consumer_test;

CREATE TABLE consumer_test.metrics_gf_test_local
(
    `metricName` LowCardinality(String),
    `serviceName` LowCardinality(String),
    `serviceRegion` LowCardinality(String),
    `serviceAz` LowCardinality(String),
    `labelKeys` Array(LowCardinality(String)),
    `labelValues` Array(LowCardinality(String)),
    `timestamp` Int64,
    `value` Float64
)
ENGINE = MergeTree()
PARTITION BY toYYYYMMDD(toDateTime(timestamp))
ORDER BY (serviceName, metricName, timestamp)
SETTINGS index_granularity = 8192;