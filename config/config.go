package config

type Shard struct {
	Name string
	Idx  int //unique
}

// describes sharding config
type Config struct {
	Shard []Shard
}
