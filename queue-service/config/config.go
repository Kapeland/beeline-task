package config

import (
	"flag"
	"time"
)

type Config struct {
	BrokerAddress string
	ReqTimeout    time.Duration
	MaxQueueSize  int
	MaxQueues     int
}

func ParseFlags() *Config {
	cfg := Config{}

	flag.StringVar(&cfg.BrokerAddress, "addr", "127.0.0.1:8080", "broker address")
	flag.DurationVar(&cfg.ReqTimeout, "t", time.Second*5, "request read timeout")
	flag.IntVar(&cfg.MaxQueueSize, "qs", 5, "queue size")
	flag.IntVar(&cfg.MaxQueues, "q", 2, "maximum number of queues")
	flag.Parse()

	return &cfg
}
