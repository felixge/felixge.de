package config

import "os"

func EnvConfig() *Config {
	c := &Config{
		Addr:    "localhost:8080",
		Service: "felixge.de",
		Env:     "dev",
		Version: "n/a",
		Datadog: true,
		EmbedFS: true,
	}
	if v := os.Getenv("PORT"); v != "" {
		c.Addr = ":" + v
	}
	if v := os.Getenv("EMBEDFS"); v == "false" {
		c.EmbedFS = false
	}
	if v := os.Getenv("DATADOG"); v == "false" {
		c.Datadog = false
	}
	if c.Datadog {
		if v := os.Getenv("DD_SERVICE"); v != "" {
			c.Service = v
		}
		if v := os.Getenv("DD_ENV"); v != "" {
			c.Env = v
		}
		if v := os.Getenv("DD_VERSION"); v != "" {
			c.Version = v
		}
	}
	return c
}

type Config struct {
	Service string
	Env     string
	Version string
	Addr    string
	Datadog bool
	EmbedFS bool
}
