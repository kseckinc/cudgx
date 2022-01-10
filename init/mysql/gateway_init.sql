use cudgx;
DROP TABLE IF EXISTS `rules`;
CREATE TABLE `rules` (
                         `id` int(11) NOT NULL AUTO_INCREMENT,
                         `metric_name` varchar(256) DEFAULT NULL,
                         `service_name` varchar(256) DEFAULT NULL,
                         `aggregate` TEXT DEFAULT NULL,
                         `filters` TEXT DEFAULT NULL,
                         `groups` TEXT DEFAULT NULL,
                         `benchmark` double DEFAULT NULL,
                         `ts` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
                         PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 ;

CREATE UNIQUE INDEX uniq_service_metric on  rules(`service_name`,`metric_name`);
