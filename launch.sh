#!/bin/bash

#The script will stop on errors
set -e

#when the script receives SIGINT (Ctrl+C), it will run the command 'killall DistributedKV' to terminate all processes called DistributedKV
trap 'killall DistributedKV' SIGINT

killall DistributedKV || true
sleep 0.1

go install .

sudo ifconfig lo0 alias 127.0.0.1 up
# for i in {2..55}; do
#   sudo ifconfig lo0 alias 127.0.0.$i up
# done


DistributedKV --db-location="name.db" --http-address=127.0.0.1:8081 --config-file="sharding.toml" --shard="name" &
DistributedKV --db-location="place.db" --http-address=127.0.0.1:8082 --config-file="sharding.toml" --shard="place" &
DistributedKV --db-location="animal.db" --http-address=127.0.0.1:8083 --config-file="sharding.toml" --shard="animal" &
DistributedKV --db-location="thing.db" --http-address=127.0.0.1:8084 --config-file="sharding.toml" --shard="thing" &

wait