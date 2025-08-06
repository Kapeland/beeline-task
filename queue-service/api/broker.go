package main

import (
	"net/http"
	iQ "queue-broker/api/adapters/inmemoryQueue"
	"queue-broker/api/adapters/rest"
	"queue-broker/config"
)

func createBrokerServer(cfg *config.Config) *http.Server {
	qBroker := iQ.NewInMemoryQueueRepository(cfg.MaxQueueSize, cfg.MaxQueues)

	mux := http.NewServeMux()

	mux.Handle("PUT /queue/{queue}", rest.PutMsg(qBroker))
	mux.Handle("GET /queue/{queue}", rest.GetMsg(qBroker))

	brokerServer := &http.Server{
		Addr:        cfg.BrokerAddress,
		ReadTimeout: cfg.ReqTimeout,
		Handler:     mux,
	}

	return brokerServer
}
