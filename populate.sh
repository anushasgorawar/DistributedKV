#!/bin/bash

for shard in {localhost:8081,localhost:8082}; do
	for i in {1..100}; do
		curl http://${shard}/set?key=key-$i\&value=value-$i
	done
done