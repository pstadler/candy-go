package main

import (
	"encoding/json"
	"fmt"
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
	HTTP_Bind string
	App       struct {
		Host string
		Port int
	}
}

var (
	json_config *Configuration
)

func config(key string) string {
	// lazy-load json file
	if json_config == nil {
		file, _ := os.Open("config.json")
		json_config = &Configuration{}
		json.NewDecoder(file).Decode(json_config)
	}

	switch key {
	case "Core":
		val := os.Getenv("CANDY_CORE")
		if val != "" {
			return val
		}
		str, _ := json.Marshal(json_config.Candy.Core)
		return string(str)
	case "View":
		val := os.Getenv("CANDY_VIEW")
		if val != "" {
			return val
		}
		str, _ := json.Marshal(json_config.Candy.View)
		return string(str)
	case "Connect":
		val := os.Getenv("CANDY_CONNECT")
		if val == "" {
			b, _ := json.Marshal(json_config.Candy.Connect)
			val = string(b)
		}
		if val != "null" {
			val = val[1 : len(val)-1]
		}
		return val
	case "HTTP_Bind":
		val := os.Getenv("CANDY_HTTP_BIND")
		if val != "" {
			return val
		}
		return json_config.HTTP_Bind
	case "App":
		val := os.Getenv("PORT")
		if val != "" {
			return fmt.Sprintf(":%s", val)
		}
		return fmt.Sprintf("%s:%d", json_config.App.Host, json_config.App.Port)
	}
	return ""
}
