package config

import (
	"fmt"
	"hash/fnv"
	"log"
)

type Shard struct {
	Name string
	Idx  int //unique
	Addr string
}

// describes sharding config
type Config struct {
	Shards []Shard
}

type Shards struct {
	Count   int
	CurrInd int
	Addrs   map[int]string
}

func ParseShards(shards []Shard, currShardName string) (*Shards, error) {
	shardCount := len(shards)
	var shardIndex int = -1
	addrs := make(map[int]string)

	for _, s := range shards {
		if _, ok := addrs[s.Idx]; ok {
			return nil, fmt.Errorf("duplicate shard index: %v", s.Idx)
		}
		addrs[s.Idx] = s.Addr
		if s.Name == currShardName {
			shardIndex = s.Idx
		}
	}
	for i := 0; i < shardCount; i++ {
		if _, ok := addrs[i]; !ok {
			return nil, fmt.Errorf("missing shard index: %v", i)
		}
	}
	if shardIndex == -1 {
		log.Fatalf("Shard %q not found.", currShardName)
	}

	log.Printf("Total shards: %v, current shard: %v, index: %v", shardCount, currShardName, shardIndex)
	return &Shards{
		Count:   shardCount,
		CurrInd: shardIndex,
		Addrs:   addrs,
	}, nil
}

func (s *Shards) GetShard(key string) int {
	hash := fnv.New64()
	hash.Write([]byte(key))
	return int(hash.Sum64() % uint64(s.Count))
}
