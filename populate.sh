#!/bin/bash

for shard in {localhost:8081,localhost:8082}; do
	for i in {1..10}; do
		curl http://${shard}/set?key=key-$RANDOM\&value=value-$RANDOM
	done
done