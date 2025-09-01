package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/anushasgorawar/DistributedKV/db"
)

// Serve contains http method handlers to be used for the db
type Server struct {
	db *db.Database
}

// newServer creates a new server instance with hundlers
func NewServer(db *db.Database) *Server {
	return &Server{db: db}
}

// Handles "get" endpoint
func (d *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Get function is called")
	r.ParseForm()
	key := r.Form.Get("key")
	value, err := d.db.GetKey(key)
	fmt.Fprintf(w, "value=%q, error:%v", value, err)
}

// Handles "set" endpoint
func (d *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Set function is called")
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")
	err := d.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "value=%q, error:%v", value, err)
}

