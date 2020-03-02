package roboot

import (
	"encoding/json"
	"log"
	"net/http"
)

// https://github.com/coreos/airlock/pull/1/files

type locker interface {
	claim(node string) (bool, error)
	release(node string) error
}

type scheduler interface {
	drain(node string) error
	uncordon(node string) error
}

type healthChecker interface {
	healthy() (bool, error)
}

type request struct {
	ClientParams struct {
		Group    string
		NodeUUID string `json:"node_uuid"`
	} `json:"client_params"`
}

func (s *server) handleSteadyState() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			log.Printf("decode steady state request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		nodeID := req.ClientParams.NodeUUID

		// uncordon node
		err := s.scheduler.uncordon(nodeID)
		if err != nil {
			log.Printf("uncordon node: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// relase lock
		err = s.locker.release(nodeID)
		if err != nil {
			log.Printf("release lock: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

func (s *server) handlePreReboot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			log.Printf("decode pre reboot request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		nodeID := req.ClientParams.NodeUUID

		// check cluster is healthy
		ok, err := s.healthChecker.healthy()
		if err != nil {
			log.Printf("cluster health check: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !ok {
			w.WriteHeader(http.StatusConflict)
			return
		}

		// aquire lock
		ok, err = s.locker.claim(nodeID)
		if err != nil {
			log.Printf("claim lock: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !ok {
			w.WriteHeader(http.StatusConflict)
			return
		}

		// drain node
		err = s.scheduler.drain(nodeID)
		if err != nil {
			log.Printf("drain node: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
