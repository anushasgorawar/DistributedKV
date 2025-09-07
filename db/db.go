package db

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

var (
	defaultBucket = []byte("default")
)

// A bolt database
type Database struct {
	db *bolt.DB
}

// a function to create a bucket in BoltDB
func (d *Database) createDefaultBucket() error {
	return d.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(defaultBucket)
		return err
	})
}

// New database instance
func NewDatabase(dbLocation string) (db *Database, closeFunc func() error, err error) {
	//open a db
	boltDB, err := bolt.Open(dbLocation, 0600, nil)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	// create a db struct
	db = &Database{db: boltDB}
	closeFunc = boltDB.Close

	//create a bucket
	if err := db.createDefaultBucket(); err != nil {
		closeFunc()
		return nil, nil, fmt.Errorf("createDefaultBucket() failed")
	}
	return db, closeFunc, nil
}

// SetKey sets a key to a value, else, returns an error
func (d *Database) SetKey(key string, value []byte) error {
	return d.db.Update(func(tx *bolt.Tx) error {

		bucket := tx.Bucket(defaultBucket)
		return bucket.Put([]byte(key), value)
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
