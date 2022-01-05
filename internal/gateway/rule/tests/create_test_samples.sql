create database cudgx_test;

use cudgx_test;

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
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 ;

CREATE UNIQUE INDEX uniq_service_metric on  rules(`service_name`,`metric_name`);


insert into rules (metric_name,service_name,aggregate,filters,benchmark,ts) values('latency','gf.cudgx.sample.pi','{
        "operation": "section_factor",
        "param": "{\\"sections\\":[10,30,40,50,100],\\"factors\\":[0.01,0.1,0.3,0.5,1,10]}"}',
        '[{"key":"serviceName","value":"gf","action":"like"},{"key":"metricName","value":"latency","action":"equals"}]',500,now()) ;


