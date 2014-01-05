package main

import (
	"encoding/json"
	"os"
)

// TODO: make this more flexible
// http://blog.golang.org/json-and-go

type Configuration struct {
	Candy struct {
		Core    interface{}
		View    interface{}
		Connect interface{}
	}
	HTTP_Bind struct {
		Host string
		Port int
		Path string
	}
	App struct {
		Host string
		Port int
	}
}

func loadConfig(path string) *Configuration {
	config, _ := os.Open(path)
	configuration := &Configuration{}
	json.NewDecoder(config).Decode(configuration)
	return configuration
}
