package db_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/anushasgorawar/DistributedKV/db"
)

func TestSetGet(t *testing.T) {
	dbName := "test.db"
	defer os.Remove(dbName)
	newDB, closeFunc, err := db.NewDatabase(dbName)
	if err != nil {
		t.Errorf("could not create db: %q", err.Error())
	}
	defer closeFunc()
	err = newDB.SetKey("table", []byte("Round"))
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

func getKey(t *testing.T, db *db.Database, key string) string {
	t.Helper()
	value, err := db.GetKey(key)
	if err != nil {
		t.Errorf("Could not set key")
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
	dbName := "test.db"
	defer os.Remove(dbName)
	newDB, closeFunc, err := db.NewDatabase(dbName)
	if err != nil {
		t.Errorf("could not create db: %q", err.Error())
	}
	defer closeFunc()
	setKey(t, newDB, "table", "Round")
	setKey(t, newDB, "monitor", "rectangle")

	value := getKey(t, newDB, "table")
	if err != nil {
		t.Errorf("could not find key: %q", err)
	}
	if value != "Round" {
		t.Errorf("Could not get key")
	}

	if err = newDB.Purge(func(key string) bool { return key == "table" }); err != nil {
		t.Errorf("Could not delete keys")
	}
	if value = getKey(t, newDB, "Round"); value != "" {
		t.Errorf("Deleted everything")
	}
	if value = getKey(t, newDB, "table"); value != "" {
		t.Errorf("Failed to delete keys")
	}
}
