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
