package main

import (
	"flag"
	"github.com/BurntSushi/toml"
	"log"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "conf", "./proxy.toml", "path to configs file")
}

func main() {
	config := DefaultConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	if config.ListenAddr == "" {
		config.ListenAddr = ":8888"
	}
	log.Fatal(RunProxy(config))
}
