package db_test

import (
	"bytes"
	"os"
	"reflect"
	"testing"

	"github.com/anushasgorawar/DistributedKV/db"
)

func createDatabase(t *testing.T, readOnly bool) *db.Database {
	t.Helper()
	dbName := "test.db"
	defer os.Remove(dbName)
	newDB, closeFunc, err := db.NewDatabase(dbName, readOnly)
	if err != nil {
		t.Errorf("could not create db: %q", err.Error())
	}
	// defer closeFunc()
	t.Cleanup(func() { closeFunc() })
	return newDB
}
func TestSetGet(t *testing.T) {
	newDB := createDatabase(t, false)
	err := newDB.SetKey("table", []byte("Round"))
	if err != nil {
		t.Errorf("could not set key: %q", err.Error())
	}
	recieved, err := newDB.GetKey("table")
	if err != nil {
		t.Errorf("could not find key: %q", err.Error())
	}
	expected := []byte("Round")

	if !reflect.DeepEqual(recieved, expected) {
		t.Errorf("expected: %q, recieved: %q", expected, recieved)
	}
}

func TestReplicaSet(t *testing.T) {
	newDB := createDatabase(t, true)
	recieved := newDB.SetKey("table", []byte("Round"))
	if recieved == nil {
		t.Errorf("could set key in replica shard: %q", recieved)
	}
}

func getKey(t *testing.T, db *db.Database, key string) string {
	t.Helper()
	value, err := db.GetKey(key)
	if err != nil {
		t.Errorf("Could not get key")
	}
	return string(value)
}
func setKey(t *testing.T, db *db.Database, key, value string) {
	t.Helper()
	if err := db.SetKey(key, []byte(value)); err != nil {
		t.Errorf("Could not set key")
	}
}
func TestPurge(t *testing.T) {
	newDB := createDatabase(t, false)
	setKey(t, newDB, "table", "Round")
	setKey(t, newDB, "monitor", "rectangle")

	value := getKey(t, newDB, "table")
	if value != "Round" {
		t.Errorf("Could not get key")
	}

	if err := newDB.Purge(func(key string) bool { return key == "table" }); err != nil {
		t.Errorf("Could not delete keys")
	}
	if value = getKey(t, newDB, "Round"); value != "" {
		t.Errorf("Deleted everything")
	}
	if value = getKey(t, newDB, "table"); value != "" {
		t.Errorf("Failed to delete keys")
	}
}

func TestDeleteReplicationKey(t *testing.T) {
	db := createDatabase(t, false)

	setKey(t, db, "party", "Great")

	k, v, err := db.GetNextKeyForReplication()
	if err != nil {
		t.Fatalf(`Unexpected error for GetNextKeyForReplication(): %v`, err)
	}

	if !bytes.Equal(k, []byte("party")) || !bytes.Equal(v, []byte("Great")) {
		t.Errorf(`GetNextKeyForReplication(): got %q, %q; want %q, %q`, k, v, "party", "Great")
	}

	if err := db.DeleteReplicatedKey([]byte("party"), []byte("Bad")); err == nil {
		t.Fatalf(`DeleteReplicationKey("party", "Bad"): got nil error, want non-nil error`)
	}

	if err := db.DeleteReplicatedKey([]byte("party"), []byte("Great")); err != nil {
		t.Fatalf(`DeleteReplicationKey("party", "Great"): got %q, want nil error`, err)
	}

	k, v, err = db.GetNextKeyForReplication()
	if err != nil {
		t.Fatalf(`Unexpected error for GetNextKeyForReplication(): %v`, err)
	}

	if k != nil || v != nil {
		t.Errorf(`GetNextKeyForReplication(): got %v, %v; want nil, nil`, k, v)
	}
}
