package camp

// https://www.youtube.com/watch?v=rWBSMsLG8po

import (
	"encoding/json"
	"log"
	"net/http"
)

func NewServer() *server {
	s := &server{}
	s.routes()
	return s
}

type server struct {
	router http.ServeMux
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) routes() {
	//	s.router.Handle("/health", s.handleHealth())
	s.router.Handle("/v1/pre-reboot", s.onlyFleetLockProtocol(s.handlePreReboot()))
	//	s.router.Handle("/v1/steady-state", s.handleSteadyState())
}

func (s *server) onlyFleetLockProtocol(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("fleet-lock-protocol") != "true" {
			http.NotFound(w, r)
		}

		h(w, r)
	}
}

func (s *server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Printf("encode json: %v", err)
		}
	}
}

// https://github.com/coreos/airlock/pull/1/files
func (s *server) handlePreReboot() http.HandlerFunc {

	type request struct {
		ClientParams struct {
			Group    string
			NodeUUID string `json:"node_uuid"`
		} `json:"client_params"`
	}
	type response struct {
		Status string `json:"status"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := s.decode(w, r, req); err != nil {
			log.Printf("handle pre reboot: %v", err)
		}

		// check all nodes are up
		// check all pods are up
		// drain node

		s.respond(w, r, response{Status: "ok"}, http.StatusOK)
	}
}
