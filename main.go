package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/anushasgorawar/DistributedKV/db"
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

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		// log.Println("Get function is called")
		// fmt.Fprintf(w, "Get function is called")

		r.ParseForm()
		key := r.Form.Get("key")
		value, err := boltDB.GetKey(key)
		// w.Write(value)
		fmt.Fprintf(w, "value=%q, error:%v", value, err)
	})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		// log.Println("Set function is called")
		// fmt.Fprintf(w, "Set function is called")
		r.ParseForm()
		key := r.Form.Get("key")
		value := r.Form.Get("value")
		err := boltDB.SetKey(key, []byte(value))
		// w.Write(value)
		fmt.Fprintf(w, "value=%v, error:%v", value, err)
	})

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
