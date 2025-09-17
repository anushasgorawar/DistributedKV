package db

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

var (
	defaultBucket = []byte("default")
	replicaBucket = []byte("replica")
)

// A bolt database
type Database struct {
	db       *bolt.DB
	readOnly bool //true if replica
}

// a function to create a bucket in BoltDB
func (d *Database) createBuckets() error {
	return d.db.Update(func(tx *bolt.Tx) (err error) {
		if _, err := tx.CreateBucketIfNotExists(defaultBucket); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists(replicaBucket); err != nil {
			return err
		}
		return err
	})
}

// New database instance
func NewDatabase(dbLocation string, readOnly bool) (db *Database, closeFunc func() error, err error) {
	//open a db
	boltDB, err := bolt.Open(dbLocation, 0600, nil)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	// boltDB.NoSync = true // flushing data to boltdb for every write is disabled
	// create a db struct
	db = &Database{db: boltDB, readOnly: readOnly}
	closeFunc = boltDB.Close

	//create a bucket
	if err := db.createBuckets(); err != nil {
		closeFunc()
		return nil, nil, fmt.Errorf("createBuckets() failed")
	}
	return db, closeFunc, nil
}

// SetKey sets a key to a value, else, returns an error
func (d *Database) SetKey(key string, value []byte) error {
	if d.readOnly {
		return errors.New("read-only mode")
	}
	return d.db.Update(func(tx *bolt.Tx) (err error) {
		if err := tx.Bucket(defaultBucket).Put([]byte(key), value); err != nil {
			return err
		}
		return tx.Bucket(replicaBucket).Put([]byte(key), value)
	})
}

// GetKey returns the vlaue of a key
func (d *Database) GetKey(key string) ([]byte, error) {
	var res []byte
	err := d.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(defaultBucket)
		res = bucket.Get([]byte(key))
		return nil //return of this anon func
	})
	if err == nil {
		return res, nil
	}
	return nil, err

}

func (d *Database) Purge(isExtra func(string) bool) error {
	var keys []string
	err := d.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(defaultBucket)
		return bucket.ForEach(func(k, v []byte) error {
			if isExtra(string(k)) {
				keys = append(keys, string(k))
			}
			return nil
		})
	})
	if err != nil {
		log.Fatal("Could not list keys: ", err)
	}
	return d.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(defaultBucket)
		for _, key := range keys {
			if err := bucket.Delete([]byte(key)); err != nil {
				log.Fatal("Could not delete key: ", err)
				return err
			}
		}
		return nil
	})
}

func (d *Database) DeleteReplicatedKey(key, value []byte) (err error) {
	return d.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(defaultBucket)

		v := bucket.Get([]byte(key)) //returns null if not exists
		if v == nil {
			return nil
		}
		if bytes.Equal(v, value) {
			return errors.New("values do not match")
		}
		return bucket.Delete(key)
	})
}
