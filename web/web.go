package web

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/anushasgorawar/DistributedKV/config"
	"github.com/anushasgorawar/DistributedKV/db"
)

// Serve contains http method handlers to be used for the db
type Server struct {
	db     *db.Database
	shards *config.Shards
}

// newServer creates a new server instance with hundlers
func NewServer(db *db.Database, shards *config.Shards) *Server {
	return &Server{
		db:     db,
		shards: shards,
	}
}

// func (d *Server) getShard(key string) int {
// 	hash := fnv.New64()
// 	hash.Write([]byte(key))
// 	return int(hash.Sum64() % uint64(d.shardcount))
// }

func (d *Server) handle(w http.ResponseWriter, shardIdx int, r *http.Request) {
	url := "http://" + d.shards.Addrs[shardIdx] + r.RequestURI

	log.Println(url)
	log.Println("Redirecting from shard:", d.shards.CurrInd, "to", shardIdx)

	res, err := http.Get(url)

	if err != nil {
		w.WriteHeader(500)
		log.Fatal("Error redirecting the request:", err)
		return
	}
	defer res.Body.Close()
	io.Copy(w, res.Body) //explain
}

// Handles "get" endpoint
func (d *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get function is called")
	r.ParseForm()
	key := r.Form.Get("key")

	shardIdx := d.shards.GetShard(key)

	// fmt.Fprintf(w, "shard=%v currShard=%v addr=%v value=%q, error:%v", shardIdx, d.shardIndex, d.addr[shardIdx], value, err)

	if shardIdx != d.shards.CurrInd {
		d.handle(w, shardIdx, r)
		return
	} else {
		value, err := d.db.GetKey(key)
		fmt.Fprintf(w, "shard=%v addr=%v value=%q, error:%v", shardIdx, d.shards.Addrs[shardIdx], value, err)
	}
}

// Handles "set" endpoint
func (d *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Set function is called")
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	shardIdx := d.shards.GetShard(key)

	if shardIdx != d.shards.CurrInd {
		d.handle(w, shardIdx, r)
	} else {
		err := d.db.SetKey(key, []byte(value))

		fmt.Fprintf(w, "shard=%v addr=%v value=%q, error:%v", shardIdx, d.shards.Addrs[shardIdx], value, err)
	}
}
