package main

import (
	"fmt"
	"github.com/enpsl/conf-reload"
	"os"
	"time"
)

type Http struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func main() {
	f := "example.toml"
	if _, err := os.Stat(f); err != nil {
		f = "_example/example.toml"
	}
	conf_reload.LoadEngine(f, conf_reload.WithLevelSplit("."), conf_reload.WithLogLevel(0))
	var http = &Http{}
	for {
		err := conf_reload.DecodeToStruct("server.http", http)
		if err != nil {
			panic(err)
		}
		fmt.Println(http)
		time.Sleep(2 * time.Second)
	}
}
