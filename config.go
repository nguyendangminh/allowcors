package main

type Config struct {
	Port     int `yaml:"port"`
	Backends []struct {
		Name   string `yaml:"name"`
		Path   string `yaml:"path"`
		Target string `yaml:"target"`
	} `yaml:"backends"`
}
