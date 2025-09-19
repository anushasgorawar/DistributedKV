package replication

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/anushasgorawar/DistributedKV/db"
)

type NextKeyValue struct {
	Key   string
	Value string
	Err   error
}

type client struct {
	db         *db.Database
	leaderAddr string
}

func Clientloop(db *db.Database, leaderAddr string) {
	// ClientLoop continuously downloads new keys from the master and applies them. Every second.
	c := &client{
		db:         db,
		leaderAddr: leaderAddr,
	}
	for {
		present, err := c.loop()
		if err != nil {
			log.Fatal("could not retrieve the next key for replication.")
			time.Sleep(time.Second) //fixme: change to time.Second after testing
			continue
		}
		if !present {
			time.Sleep(time.Second) //fixme: change to time.Second after testing
		}
	}

}

func (c *client) loop() (present bool, err error) {
	resp, err := http.Get("http://" + c.leaderAddr + "/next-replication-key")
	if err != nil {
		return false, err
	}
	var res NextKeyValue
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return false, err
	}
	if res.Err != nil {
		return false, res.Err
	}
	if res.Key == "" {
		return false, nil
	}

	if err := c.db.ReplicaSetKey(res.Key, []byte(res.Value)); err != nil {
		return false, err
	}

	if err := c.deleteFromReplicationQueue(res.Key, res.Value); err != nil {
		log.Fatalf("could not delete key: %v", res.Key)
	}

	defer resp.Body.Close()
	log.Printf("Next ky-value: %+v", res)
	return true, nil
}

func (c *client) deleteFromReplicationQueue(key, value string) error {
	keyvalue := url.Values{}
	keyvalue.Set("key", key)
	keyvalue.Set("value", value)
	log.Printf("Deleting key=%v,value=%v on replication queue on %q", key, value, c.leaderAddr)
	url := "http://" + c.leaderAddr + "/delete-next-replication-key?" + keyvalue.Encode()
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	result, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if !bytes.Equal(result, []byte("ok")) {
		return errors.New(string(result))
	}
	return nil
}
