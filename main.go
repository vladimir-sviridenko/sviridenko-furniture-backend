package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

func getConfig() (*Config, error) {
	file, ok := os.LookupEnv("CONFIG_FILE")
	if !ok {
		file = "config.json"
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	cfg, err := getConfig()
	if err != nil {
		panic(err)
	}
	srv := http.Server{
		Addr: cfg.Server.Addr,
	}
	srv.ListenAndServe()
}
