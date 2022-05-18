#!/bin/bash
set -e

trap 'killall kv-go' SIGINT

cd $(dirname $0)

killall kv-go || true
sleep 0.5

go install -v

kv-go  --db-location=$PWD/bkk.db --http-addr=127.0.0.1:8080 --config-file=$PWD/sharding.toml --shard=BKK &
kv-go  --db-location=$PWD/cnx.db --http-addr=127.0.0.1:8081 --config-file=$PWD/sharding.toml --shard=CNX &
kv-go  --db-location=$PWD/hkt.db --http-addr=127.0.0.1:8082 --config-file=$PWD/sharding.toml --shard=HKT &

wait