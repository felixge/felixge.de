package config

import "os"

func EnvConfig() *Config {
	c := &Config{
		Addr: ":" + os.Getenv("PORT"),
	}
	if c.Addr == ":" {
		c.Addr = "localhost:8080"
	}
	return c
}

type Config struct {
	Addr string
}
