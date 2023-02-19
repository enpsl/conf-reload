package conf_reload

import (
	"os"
	"testing"
	"time"
)

type Http struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func init() {
	f := "example.toml"
	if _, err := os.Stat(f); err != nil {
		f = "_example/example.toml"
	}
	LoadEngine(f, WithLevelSplit("."), WithLogLevel(1))
}

func TestDecodeToStruct(t *testing.T) {
	var http = &Http{}
	err := DecodeToStruct("server.http", http)
	if err != nil {
		t.Fatal(err)
	}
	if http.Port != 8080 || http.Host != "0.0.0.0" {
		t.Errorf("host=%s, want=%s| port=%d, want=%d", http.Host, "0.0.0.0", 8080, http.Port)
	}
}

func TestGetInt(t *testing.T) {
	i := GetInt("server.http.port")
	if i != 8080 {
		t.Errorf("god=%d, want=%d", 8080, i)
	}
}

func TestGetString(t *testing.T) {
	i := GetString("server.http.host")
	if i != "0.0.0.0" {
		t.Errorf("god=%s, want=%s", "0.0.0.0", i)
	}
}

func TestGetBool(t *testing.T) {
	i := GetBool("server.config.connection")
	if i != false {
		t.Errorf("god=%t, want=%t", i, false)
	}
}

func TestGetDuration(t *testing.T) {
	i := GetDuration("server.config.timeout")
	if i != 10*time.Second {
		t.Errorf("god=%d, want=%d", i, 10*time.Second)
	}
}

func TestGetTime(t *testing.T) {
	i := GetTime("server.config.publish")
	parse, _ := time.Parse("2006-01-02", "2023-02-19")
	if i.Unix() != parse.Unix() {
		t.Errorf("god=%d, want=%d", i.Unix(), parse.Unix())
	}
}

func TestGetSlice(t *testing.T) {
	i := GetSlice("server.config.depends")
	for _, v := range i {
		switch v {
		case "tcp":
		case "ip":
		default:
			t.Errorf("god=%s, want=%s", v, "tcp or ip")
		}
	}
}
