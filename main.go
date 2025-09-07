package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/anushasgorawar/DistributedKV/config"
	"github.com/anushasgorawar/DistributedKV/db"
	"github.com/anushasgorawar/DistributedKV/web"
)

var (
	dbLocation = flag.String("db-location", "", "the path to boltDb Location")
	httpAddr   = flag.String("http-address", "127.0.0.1:8080", "HTTP Host and Port")
	configFile = flag.String("config-file", "sharding.toml", "config file for shards")
	shard      = flag.String("shard", "", "The name of the shard for they key-value")
)

func parseFlags() {
	flag.Parse()
	if *dbLocation == "" {
		log.Fatal("Must provide db-location.")
	}
}

func main() {
	parseFlags()

	var c config.Config
	// log.Printf("%#v", &c)

	if _, err := toml.DecodeFile(*configFile, &c); err != nil {
		log.Fatalf("Unable to decode config File %v, error: %v", *configFile, err)
	}
	//to see if config is parsed properly
	// log.Printf("%#v", &c)

	shards, err := config.ParseShards(c.Shards, *shard)
	if err != nil {
		log.Fatalf("Unable to parse shards: %v", err)
	}

	log.Printf("%#v", &shards)

	// var shardCount int
	// var shardIndex int = -1
	// addrs := make(map[int]string)
	// shardCount = len(c.Shards)
	// for _, s := range c.Shards {
	// 	addrs[s.Idx] = s.Addr
	// 	if s.Idx+1 > shardCount {
	// 		shardCount = s.Idx + 1
	// 	}
	// 	if s.Name == *shard {
	// 		shardIndex = s.Idx
	// 	}
	// }
	// if shardIndex == -1 {
	// 	log.Fatal("Shard not found.")
	// }

	// log.Printf("Total shards: %v, current shard: %v, index: %v", shardCount, *shard, shardIndex)

	boltDB, closefunc, err := db.NewDatabase(*dbLocation)
	defer closefunc()
	if err != nil {
		log.Fatalf("NewDatabase(%q):%v", *dbLocation, err) //if dbLocation is unclear
	}

	server := web.NewServer(boltDB, shards)
	http.HandleFunc("/get", server.GetHandler)

	http.HandleFunc("/set", server.SetHandler)

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
