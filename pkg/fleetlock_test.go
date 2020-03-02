package roboot

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

type lockMock struct{}

func (l lockMock) claim(node string) (bool, error) {
	return true, nil
}

func (l lockMock) release(node string) error {
	return nil
}

type schedulerMock struct{}

func (l schedulerMock) drain(node string) error {
	return nil
}

func (l schedulerMock) uncordon(node string) error {
	return nil
}

type healthCheckerMock struct{}

func (h healthCheckerMock) healthy() (bool, error) {
	return true, nil
}

func Test_server_handleSteadyState(t *testing.T) {
	is := is.New(t)
	srv := NewServer()
	srv.locker = lockMock{}
	srv.scheduler = schedulerMock{}

	p := request{
		ClientParams: struct {
			Group    string
			NodeUUID string `json:"node_uuid"`
		}{
			Group:    "default",
			NodeUUID: "node1",
		},
	}
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(p)
	is.NoErr(err)

	r := httptest.NewRequest(http.MethodGet, "/", &buf)
	w := httptest.NewRecorder()

	srv.handleSteadyState()(w, r)

	defer w.Result().Body.Close()
	is.Equal(w.Result().StatusCode, http.StatusOK)
}

func Test_server_handleSteadyState_no_body(t *testing.T) {
	is := is.New(t)
	srv := NewServer()

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	srv.handleSteadyState()(w, r)

	is.Equal(w.Result().StatusCode, http.StatusBadRequest)
}
