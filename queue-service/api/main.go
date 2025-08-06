package main

import "queue-broker/config"

func main() {
	cfg := config.ParseFlags()

	if err := NewService().Start(cfg); err != nil {
		//slog.Error(err.Error(), "description", "Failed creating service")
		return
	}
}
