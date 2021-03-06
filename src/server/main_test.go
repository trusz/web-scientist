package server_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/sprinteins/web-scientist/server"
	. "github.com/sprinteins/web-scientist/server/test_helpers"
)

const PROTOCOL = "http"
const HOST = "localhost"
const PORT = "2345"

var scientist *server.Server

func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	teardown()
	os.Exit(retCode)
}

func setup() {
	scientist = server.New(HOST, PORT)
	go scientist.Start()

}

func teardown() {
	scientist.Stop()
}

func Test_By_Failed_Experiment_Reference_Sent(t *testing.T) {

	var reference, experiment = CreateNonEqualMocks()

	scientist.SetReference(reference.Address())
	scientist.SetExperiment(experiment.Address())

	var message = "TeSt"
	var payload = []byte(message)

	var resp, err = http.Post(scientist.Address(), "text/plain", bytes.NewBuffer(payload))
	Ok(t, err)

	var header = resp.Header.Get("X-web-scientist-type")
	Equals(t, "reference", header)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	Ok(t, err)

	var respPayload = string(body)
	Equals(t, respPayload, message)

	reference.Stop()
	experiment.Stop()

}

func Test_By_Successfull_Experiment_Experiment_Sent(t *testing.T) {

	var reference, experiment = CreateEqualMocks()

	scientist.SetReference(reference.Address())
	scientist.SetExperiment(experiment.Address())

	var message = "TeSt"
	var payload = []byte(message)

	var resp, err = http.Post(scientist.Address(), "text/plain", bytes.NewBuffer(payload))
	Ok(t, err)

	var header = resp.Header.Get("X-web-scientist-type")
	Equals(t, "experiment", header)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	Ok(t, err)

	var respPayload = string(body)
	Equals(t, respPayload, message)

	reference.Stop()
	experiment.Stop()
}
