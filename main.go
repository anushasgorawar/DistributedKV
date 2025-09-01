package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/anushasgorawar/DistributedKV/db"
	"github.com/anushasgorawar/DistributedKV/web"
)

var (
	dbLocation = flag.String("db-location", "", "the path to boltDb Location")
	httpAddr   = flag.String("http-address", "127.0.0.1:8080", "HTTP Host and Port")
)

func parseFlags() {
	flag.Parse()
	if *dbLocation == "" {
		log.Fatal("Must provide db location.")
	}
}

func main() {
	parseFlags()
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
