package web

import (
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"

	"github.com/anushasgorawar/DistributedKV/db"
)

// Serve contains http method handlers to be used for the db
type Server struct {
	db         *db.Database
	shardIndex int
	shardcount int
	addr       map[int]string
}

// newServer creates a new server instance with hundlers
func NewServer(db *db.Database, shardCount, shardIndex int, addr map[int]string) *Server {
	return &Server{
		db:         db,
		shardIndex: shardIndex,
		shardcount: shardCount,
		addr:       addr,
	}
}
func (d *Server) getShard(key string) int {
	hash := fnv.New64()
	hash.Write([]byte(key))
	return int(hash.Sum64() % uint64(d.shardcount))
}

func (d *Server) handle(w http.ResponseWriter, shardIdx int, requestURI string) {

	res, err := http.Get("http://" + d.addr[shardIdx] + requestURI)
	log.Println("http://" + d.addr[shardIdx] + requestURI)
	log.Println("Redirecting from shard:", d.shardIndex, "to", shardIdx)
	if err != nil {
		log.Fatal("URL not found")
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

	shardIdx := d.getShard(key)
	value, err := d.db.GetKey(key)

	// fmt.Fprintf(w, "shard=%v currShard=%v addr=%v value=%q, error:%v", shardIdx, d.shardIndex, d.addr[shardIdx], value, err)

	if shardIdx != d.shardIndex {
		d.handle(w, shardIdx, r.RequestURI)
	}
	fmt.Fprintf(w, "shard=%v addr=%v value=%q, error:%v", shardIdx, d.addr[shardIdx], value, err)

}

// Handles "set" endpoint
func (d *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Set function is called")
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	shardIdx := d.getShard(key)
	err := d.db.SetKey(key, []byte(value))

	if shardIdx != d.shardIndex {
		d.handle(w, shardIdx, r.RequestURI)
	}
	fmt.Fprintf(w, "shard=%v addr=%v value=%q, error:%v", shardIdx, d.addr[shardIdx], value, err)
}
