package config_test

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/anushasgorawar/DistributedKV/config"
)

// func TestMain(m *testing.M) {

// }

func TestConfigparse(t *testing.T) {
	contents := `[[shards]]
	name="proper"
	idx=0
	addr="127.0.0.1:8086"
	`
	expected := config.Config{
		Shards: []config.Shard{
			{
				Name: "proper",
				Idx:  0,
				Addr: "127.0.0.1:8086",
			},
		},
	}
	filename := "recieved.toml"
	_, err := os.Create(filename)
	defer os.Remove(filename)
	if err != nil {
		log.Fatal("Could not create file.")
	}
	err = os.WriteFile(filename, []byte(contents), 0644)
	if err != nil {
		log.Fatal("Could not write config into file")
	}
	recieved, err := config.ParseConfig(filename)
	if err != nil {
		log.Fatal("Could not Parse config.")
	}
	if !reflect.DeepEqual(recieved, expected) {
		t.Errorf("expected: %q, recieved: %q", expected, recieved)
	}
}
