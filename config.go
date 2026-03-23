package main

import (
    "encoding/json"
	"os"
)

type Config struct {
    SaveFile string `json:"save_file"`
    User     string `json:"user"`
    Pass     string `json:"pass"`
    Upload   string `json:"upload"`
	Listen   string `json:"listen"`
	Crt   string `json:"crt"`
	Key   string `json:"key"`
}

func resolv(config string) (Config,error) {
	var cfg Config
	file, err := os.ReadFile(config)
	if err != nil {
    return cfg,err
	}
	err = json.Unmarshal(file, &cfg)
	if err != nil {
    return cfg,err
	}
	return cfg,nil
}