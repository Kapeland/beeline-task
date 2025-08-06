package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"queue-broker/api/core"
	"strconv"
)

const maxMsgLen = 1048576 // 1 mb

func PutMsg(qRepo core.QueueRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		qName := r.PathValue("queue")
		if qName == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// По-хорошему, думаю стоит ограничивать длину входного сообщения.
		// Либо на уровне самого запроса, либо уже при попытке положить в очередь.
		bodyReader := http.MaxBytesReader(w, r.Body, maxMsgLen)
		defer bodyReader.Close()

		var req PutMsgRequest
		decoder := json.NewDecoder(bodyReader)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if req.Msg == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := qRepo.Put(qName, core.Message(*req.Msg)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func GetMsg(qRepo core.QueueRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		qName := r.PathValue("queue")
		//По-хорошему нужно ограничить длину названия очереди
		if qName == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// По-хорошему нужно ограничить максимальный таймаут, чтобы не задудосили
		timeoutReq := r.URL.Query().Get("timeout")
		var timeout int
		if timeoutReq != "" {
			var err error
			timeout, err = strconv.Atoi(timeoutReq)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		qMsg, err := qRepo.Get(qName, timeout)
		if err != nil {
			switch {
			case errors.Is(err, core.ErrQueueGetTimeout),
				errors.Is(err, core.ErrQueueNotFound):
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		resp := &GetMsgResponse{
			Msg: string(qMsg),
		}

		if err := json.NewEncoder(w).Encode(resp); err != nil {
			return
		}
	}
}

// По логике вещей, нужно предусмотреть возможность удаления очереди из такого репозитория
// На данный момент, если очереди заполнятся, то всё.
