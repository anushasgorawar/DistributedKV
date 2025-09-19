#!/bin/bash

#The script will stop on errors
set -e

#when the script receives SIGINT (Ctrl+C), it will run the command 'killall DistributedKV' to terminate all processes called DistributedKV
trap 'killall DistributedKV' SIGINT

killall DistributedKV || true
sleep 0.1

go install .

# sudo ifconfig lo0 alias 127.0.0.1 up
for i in {2..55}; do
  sudo ifconfig lo0 alias 127.0.0.$i up
done
sudo ifconfig lo0 alias 127.0.0.2
sudo ifconfig lo0 alias 127.0.0.22
sudo ifconfig lo0 alias 127.0.0.3
sudo ifconfig lo0 alias 127.0.0.33
sudo ifconfig lo0 alias 127.0.0.4
sudo ifconfig lo0 alias 127.0.0.44
sudo ifconfig lo0 alias 127.0.0.5
sudo ifconfig lo0 alias 127.0.0.55

DistributedKV --db-location="name.db" --http-address=127.0.0.2:8080 --config-file="sharding.toml" --shard="name" &
DistributedKV --db-location="name-replica.db" --http-address=127.0.0.22:8080 --config-file="sharding.toml" --shard="name" --replica &
DistributedKV --db-location="place.db" --http-address=127.0.0.3:8080 --config-file="sharding.toml" --shard="place" &
DistributedKV --db-location="place-replica.db" --http-address=127.0.0.33:8080 --config-file="sharding.toml" --shard="place" --replica &
DistributedKV --db-location="animal.db" --http-address=127.0.0.4:8080 --config-file="sharding.toml" --shard="animal" &
DistributedKV --db-location="animal-replica.db" --http-address=127.0.0.44:8080 --config-file="sharding.toml" --shard="animal" --replica &
DistributedKV --db-location="thing-replica.db" --http-address=127.0.0.55:8080 --config-file="sharding.toml" --shard="thing" --replica &
DistributedKV --db-location="thing.db" --http-address=127.0.0.5:8080 --config-file="sharding.toml" --shard="thing" &


wait