package main

import (
	conf_relod "conf-reload"
	"fmt"
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
	conf_relod.LoadEngine(f, conf_relod.WithLevelSplit("."), conf_relod.WithLogLevel(0))
	var http = &Http{}
	for {
		err := conf_relod.DecodeToStruct("server.http", http)
		if err != nil {
			panic(err)
		}
		fmt.Println(http)
		time.Sleep(2 * time.Second)
	}
}
