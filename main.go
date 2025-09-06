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
	configFile = flag.String("configfile", "sharding.toml", "config file for shards")
)

func parseFlags() {
	flag.Parse()
	if *dbLocation == "" {
		log.Fatal("Must provide db location.")
	}
}

func main() {
	parseFlags()

	var c config.Config

	if _, err := toml.DecodeFile(*configFile, &c); err != nil {
		log.Fatalf("Unable to decode config File %v, error: %v", *configFile, err.Error())
	}

	//to see if config is parsed properly
	log.Printf("%#v", &c)

	boltDB, closefunc, err := db.NewDatabase(*dbLocation)
	defer closefunc()
	if err != nil {
		log.Fatalf("NewDatabase(%q):%v", *dbLocation, err) //if dbLocation is unclear
	}
	server := web.NewServer(boltDB)
	http.HandleFunc("/get", server.GetHandler)

	http.HandleFunc("/set", server.SetHandler)

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
