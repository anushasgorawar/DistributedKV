package config_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/anushasgorawar/DistributedKV/config"
)

func createConfig(t *testing.T, contents string) config.Config {
	t.Helper()
	filename := "recieved.toml"
	_, err := os.Create(filename)
	defer os.Remove(filename)
	if err != nil {
		t.Fatal("Could not create file.")
	}
	err = os.WriteFile(filename, []byte(contents), 0644)
	if err != nil {
		t.Fatal("Could not write config into file")
	}
	recieved, err := config.ParseConfig(filename)
	if err != nil {
		t.Fatal("Could not Parse config.")
	}
	return recieved
}
func TestParseConfig(t *testing.T) {
	recieved := createConfig(t, `[[shards]]
	name="proper"
	idx=0
	addr="127.0.0.1:8086"
	`)

	expected := config.Config{
		Shards: []config.Shard{
			{
				Name: "proper",
				Idx:  0,
				Addr: "127.0.0.1:8086",
			},
		},
	}
	if !reflect.DeepEqual(recieved, expected) {
		t.Errorf("expected: %q, recieved: %q", expected, recieved)
	}
}

func TestParseShards(t *testing.T) {
	con := createConfig(t, `
	[[shards]]
	name="adverb"
	idx=0
	addr="127.0.0.1:8086"
	[[shards]]
	name="adjective"
	idx=1
	addr="127.0.0.1:8087"
	`)

	recieved, err := config.ParseShards(con.Shards, "adverb")

	if err != nil {
		t.Errorf("Couldnot parse shards: %q", err.Error())
	}

	expected := config.Shards{
		Count:   2,
		CurrInd: 0,
		Addrs: map[int]string{
			0: "127.0.0.1:8086",
			1: "127.0.0.1:8087",
		},
	}
	if !reflect.DeepEqual(*recieved, expected) {
		t.Errorf("expected: %q, recieved: %q", expected, recieved)
	}
}
