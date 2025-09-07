package config

type Shard struct {
	Name string
	Idx  int //unique
	Addr string
}

// describes sharding config
type Config struct {
	Shards []Shard
}
