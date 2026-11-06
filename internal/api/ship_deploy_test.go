package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

// deployCtx is a context whose git config holds the deploy settings.
func deployCtx(url, timeout string, stderr *bytes.Buffer) *Ctx {
	return &Ctx{
		Stderr: stderr,
		Getenv: func(k string) string {
			if k == "RHIZA_DEPLOY_INTERVAL" {
				return "5ms"
			}
			return ""
		},
		Git: func(args ...string) (string, error) {
			if len(args) == 3 && args[0] == "config" && args[1] == "--get" {
				switch args[2] {
				case DeployURLKey:
					return url, nil
				case DeployTimeoutKey:
					return timeout, nil
				}
			}
			return "", http.ErrNotSupported
		},
	}
}

func fakeHealth(statuses ...int) (*httptest.Server, *atomic.Int32) {
	var n atomic.Int32
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i := int(n.Add(1)) - 1
		w.WriteHeader(statuses[min(i, len(statuses)-1)])
	})), &n
}

func TestShipDeployRecovers(t *testing.T) {
	srv, n := fakeHealth(503, 503, 200)
	defer srv.Close()
	r, err := runShipDeploy(deployCtx(srv.URL, "10s", &bytes.Buffer{}), ShipDeployRequest{})
	if err != nil {
		t.Fatal(err)
	}
	resp := r.(ShipDeployResponse)
	if !resp.Healthy || resp.ExitCode() != 0 || resp.Polls != 3 || n.Load() != 3 {
		t.Errorf("resp = %+v after %d requests; want healthy after 3 polls", resp, n.Load())
	}
}

func TestShipDeployTimesOut(t *testing.T) {
	srv, _ := fakeHealth(503)
	defer srv.Close()
	var stderr bytes.Buffer
	r, err := runShipDeploy(deployCtx(srv.URL, "100ms", &stderr), ShipDeployRequest{})
	if err != nil {
		t.Fatal(err)
	}
	resp := r.(ShipDeployResponse)
	if resp.Healthy || resp.ExitCode() != 1 || resp.Last != "HTTP 503" {
		t.Errorf("resp = %+v; want unhealthy with last HTTP 503", resp)
	}
	var human bytes.Buffer
	resp.Human(&human)
	if !strings.Contains(human.String(), "still answered HTTP 503") || !strings.Contains(stderr.String(), "HTTP 503") {
		t.Errorf("human = %q, stderr = %q; want the last status", human.String(), stderr.String())
	}
}

func TestShipDeployNoURLSkips(t *testing.T) {
	r, err := runShipDeploy(deployCtx("", "10s", &bytes.Buffer{}), ShipDeployRequest{})
	if err != nil {
		t.Fatal(err)
	}
	resp := r.(ShipDeployResponse)
	if !resp.Skipped || resp.Polls != 0 || resp.ExitCode() != 0 {
		t.Errorf("resp = %+v; want skipped without polling", resp)
	}
}
