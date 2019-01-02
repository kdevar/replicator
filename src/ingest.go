package main

import (
	"fmt"
	"github.com/siddontang/go-mysql/canal"
	"log"
)

type Ingest struct {
	Config  *Config
	Handler canal.EventHandler
}

func (i *Ingest) Run() {
	cfg := getConfig(i.Config)

	c, err := canal.NewCanal(cfg)

	if err != nil {
		log.Fatalf("An Error Happened\n%v", err)
	}

	coordinates, err := c.GetMasterPos()

	if err != nil {
		log.Fatalf("An Error Happened\n%v", err)
	}

	c.SetEventHandler(i.Handler)

	c.RunFrom(coordinates)
}

func getConfig(c *Config) *canal.Config {
	cfg := canal.NewDefaultConfig()
	cfg.Addr = fmt.Sprintf("%s:%d", c.Host, c.Port)
	cfg.User = c.User
	cfg.Password = c.Password
	cfg.Flavor = c.Flavor
	cfg.IncludeTableRegex = c.IncludeTables
	cfg.Dump.ExecutionPath = ""
	return cfg
}

func NewIngestServer(c *Config, h canal.EventHandler) *Ingest {
	return &Ingest{
		Config:  c,
		Handler: h,
	}
}
