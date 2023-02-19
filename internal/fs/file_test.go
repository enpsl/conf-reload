package fs

import (
	"conf-reload/internal/log"
	"os"
	"testing"
)

func TestFsBrokerParse(t *testing.T) {
	test := struct {
		path        string
		write       string
		key         string
		wantedValue int64
	}{
		path:        "./test.toml",
		write:       "port=8080",
		key:         "port",
		wantedValue: 8080,
	}
	fd, err := os.Create(test.path)
	defer func() {
		fd.Close()
		os.Remove(test.path)
	}()
	if err != nil {
		t.Fatal(err)
	}
	_, err = fd.WriteString(test.write)
	if err != nil {
		t.Fatal(err)
	}
	logger := log.NewLogger(nil)
	err, broker := NewFs(test.path, logger)
	if err != nil {
		t.Fatal(err)
	}
	content, err := broker.LoadContent()
	if err != nil {
		t.Fatal(err)
	}
	err, m := broker.Parse(content)
	if err != nil {
		t.Fatal(err)
	}
	if m["port"].(int64) != test.wantedValue {
		t.Errorf("broker.Parse(%s) outputted %d, should equal value %d",
			test.path, content, test.wantedValue)
	}
}
