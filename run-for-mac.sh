#!/usr/bin/env bash

read -r -p "Are you sure to install mysql? [Y/N] " input

case $input in
    [yY][eE][sS]|[yY])
		echo "Installing mysql..."
    # deploy mysql
    docker run -d --name bridgx_db -e MYSQL_ROOT_PASSWORD=mtQ8chN2 -e MYSQL_DATABASE=bridgx -e MYSQL_USER=gf -e MYSQL_PASSWORD=db@galaxy-future.com -p 3306:3306 -v $(pwd)/init/mysql:/docker-entrypoint-initdb.d yobasystems/alpine-mariadb:10.5.11
		;;

    [nN][oO]|[nN])
		echo "Skip mysql install, please check conf/config.yml mysql config, and import init/mysql/* to existing mysql for first install."
    ;;
    *)
		echo "Invalid input..."
		exit 1
		;;
esac

read -r -p "Are you sure to install etcd? [Y/N] " input

case $input in
    [yY][eE][sS]|[yY])
		echo "Installing etcd..."
    # deploy etcd
    docker run -d --name bridgx_etcd -e ALLOW_NONE_AUTHENTICATION=yes -e ETCD_ADVERTISE_CLIENT_URLS=http://etcd:2379 -p 2379:2379 -p 2380:2380 bitnami/etcd:3
		;;

    [nN][oO]|[nN])
		echo "Skip etcd install, please check conf/config.yml etcd config."
    ;;
    *)
		echo "Invalid input..."
		exit 1
		;;
esac

# deploy api
sed "s/127.0.0.1/host.docker.internal/g" $(pwd)/conf/config.yml.prod > $(pwd)/conf/config.yml.mac
docker run -d --name bridgx_api --add-host host.docker.internal:host-gateway -v $(pwd)/conf/config.yml.mac:/home/tiger/api/conf/config.yml.prod -p 9099:9090 galaxyfuture/bridgx-api:latest bin/wait-for-api.sh
# deploy sheduler
docker run -d --name bridgx_scheduler --add-host host.docker.internal:host-gateway -v $(pwd)/conf/config.yml.mac:/home/tiger/scheduler/conf/config.yml.prod galaxyfuture/bridgx-scheduler:latest bin/wait-for-scheduler.sh














#!/usr/bin/env bash

read -r -p "Are you sure to install mysql? [Y/N] " input

case $input in
    [yY][eE][sS]|[yY])
		echo "Installing mysql..."
    # deploy mysql
    docker run -d --name cudgx_db -e MYSQL_ROOT_PASSWORD=mtQ8chN2 -e MYSQL_DATABASE=bridgx -e MYSQL_USER=gf -e MYSQL_PASSWORD=db@galaxy-future.com -p 3336:3306 -v $(pwd)/init/mysql:/docker-entrypoint-initdb.d yobasystems/alpine-mariadb:10.5.11
		;;

    [nN][oO]|[nN])
		echo "Skip mysql install, please check conf/config.yml mysql config, and import init/mysql/* to existing mysql for first install."
    ;;
    *)
		echo "Invalid input..."
		exit 1
		;;
esac

read -r -p "Are you sure to install kafka? [Y/N] " input

case $input in
    [yY][eE][sS]|[yY])
		echo "Installing kafka..."
    # deploy kafka
    docker run -d -p 2181:2181 --name zookeeper confluentinc/cp-zookeeper:6.2.0

    docker run  -d --name kafka \
        -p 9092:9092 \
        -e KAFKA_BROKER_ID=0 \
        -e KAFKA_ZOOKEEPER_CONNECT=127.0.0.1:2181 \
        -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://127.0.0.1:9092 \
        -e KAFKA_LISTENERS=PLAINTEXT://0.0.0.0:9092 -t confluentinc/cp-kafka:6.2.0
		;;

    [nN][oO]|[nN])
		echo "Skip kafka install, please check conf/api.json kafka config."
    ;;
    *)
		echo "Invalid input..."
		exit 1
		;;
esac

read -r -p "Are you sure to install clickhouse? [Y/N] " input

case $input in
    [yY][eE][sS]|[yY])
		echo "Installing clickhouse..."
    # deploy clickhouse
    docker run -d --name clickhouse --ulimit nofile=262144:262144 yandex/clickhouse-server
    ;;
    [nN][oO]|[nN])
		echo "Skip clickhouse install, please check conf/api.json clickhouse config."
    ;;
    *)
		echo "Invalid input..."
		exit 1
		;;
esac

# deploy api
sed "s/127.0.0.1/host.docker.internal/g" $(pwd)/conf/api.json > $(pwd)/conf/api.json.mac
docker run -d --name cudgx_api --add-host host.docker.internal:host-gateway -v $(pwd)/conf/api.json:/home/tiger/api/conf/api.json galaxyfuture/cudgx-api:latest bin/wait-for-api.sh
# deploy gateway
sed "s/127.0.0.1/host.docker.internal/g" $(pwd)/conf/gateway.json > $(pwd)/conf/gateway.json.mac
docker run -d --name cudgx_gateway --add-host host.docker.internal:host-gateway -v $(pwd)/conf/gateway.json:/home/tiger/gateway/conf/gateway.json  galaxyfuture/cudgx-gateway:latest bin/wait-for-scheduler.sh
# deploy consumer
sed "s/127.0.0.1/host.docker.internal/g" $(pwd)/conf/consumer.json > $(pwd)/conf/consumer.json.mac
docker run -d --name cudgx_consumer --add-host host.docker.internal:host-gateway -v $(pwd)/conf/consumer.json:/home/tiger/api/conf/consumer.json galaxyfuture/cudgx-consumer:latest bin/wait-for-api.sh
# deploy pi
docker run -d --name cudgx_sample_pi --add-host host.docker.internal:host-gateway galaxyfuture/cudgx-sample-pi:latest bin/wait-for-api.sh
# deploy benchmark
docker run -d --name cudgx_sample_benchmark --add-host host.docker.internal:host-gateway galaxyfuture/cudgx-sample-benchmark:latest bin/wait-for-api.sh
