use cudgx;

DROP TABLE IF EXISTS `predict_rules`;
CREATE TABLE `predict_rules`
(
    `id`                 INT(11) NOT NULL AUTO_INCREMENT,
    `name`               VARCHAR(255) NOT NULL,
    `service_name`       VARCHAR(255) NOT NULL,
    `cluster_name`       VARCHAR(255) NOT NULL,
    `metric_name`        VARCHAR(255) NOT NULL,
    `benchmark_qps`      INT(11) NOT NULL,
    `min_redundancy`     INT(11) NOT NULL,
    `max_redundancy`     INT(11) NOT NULL,
    `min_instance_count` INT(11) NOT NULL,
    `max_instance_count` INT(11) NOT NULL,
    `execute_ratio`      INT(11) NOT NULL,
    `status`             VARCHAR(50)  NOT NULL DEFAULT 'enable',
    `created_time`       INT(11) NOT NULL,
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE INDEX `uniq_name` (`name`) USING BTREE,
    UNIQUE INDEX `uniq_cname_sname_mname` (`service_name`, `cluster_name`, `metric_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;