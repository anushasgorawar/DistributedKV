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
	// boltDB.NoSync = true // flushing data to boltdb for every write is disabled
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
