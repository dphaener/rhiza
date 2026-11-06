package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"
)

func init() {
	register(Operation{
		Name:     "ship.deploy",
		Summary:  "After the merge and the ship command, poll the repo's deploy health URL (rhiza.ship.deployUrl) until it answers 200 or rhiza.ship.deployTimeout passes; with no URL set, do nothing",
		Request:  reflect.TypeFor[ShipDeployRequest](),
		Response: reflect.TypeFor[ShipDeployResponse](),
		Run:      runShipDeploy,
	})
}

// Deploy check timing: how often the URL is polled (RHIZA_DEPLOY_INTERVAL
// overrides it, for tests) and how long one request may take.
const (
	deployInterval   = 10 * time.Second
	deployReqTimeout = 10 * time.Second
)

// ShipDeployRequest takes nothing: the URL and timeout are daemon settings.
type ShipDeployRequest struct{}

// ShipDeployResponse is how the deploy check went.
type ShipDeployResponse struct {
	URL     string `json:"url,omitempty" desc:"The health URL; empty if none is set, which skips the check"`
	Skipped bool   `json:"skipped" desc:"True if no deploy health URL is set"`
	Healthy bool   `json:"healthy" desc:"True once the URL answered 200"`
	Polls   int    `json:"polls" desc:"How many times the URL was asked"`
	Last    string `json:"last,omitempty" desc:"The last status it saw: an HTTP status or an error"`
	Timeout string `json:"timeout,omitempty" desc:"The deploy timeout in force"`
}

// ExitCode is 0 when the check passed or was skipped, else 1.
func (r ShipDeployResponse) ExitCode() int {
	if r.Skipped || r.Healthy {
		return 0
	}
	return 1
}

// Human is one line: how the check stands.
func (r ShipDeployResponse) Human(w io.Writer) {
	switch {
	case r.Skipped:
		fmt.Fprintln(w, "No deploy health URL set (rhiza.ship.deployUrl): skipping the deploy check.")
	case r.Healthy:
		fmt.Fprintf(w, "Deploy healthy: %s answered 200 after %d %s.\n", r.URL, r.Polls, shipPlural(r.Polls, "poll", "polls"))
	default:
		fmt.Fprintf(w, "Deploy not healthy: %s still answered %s after %d %s over %s.\n", r.URL, r.Last, r.Polls, shipPlural(r.Polls, "poll", "polls"), r.Timeout)
	}
}

func runShipDeploy(ctx *Ctx, _ any) (any, error) {
	if err := refuseShipDryRun(ctx, "ship.deploy"); err != nil {
		return nil, err
	}
	rows, err := settingsOf(ctx, nil)
	if err != nil {
		return nil, err
	}
	url := settingValue(rows, DeployURLKey)
	if url == "" {
		return ShipDeployResponse{Skipped: true}, nil
	}
	timeout, err := time.ParseDuration(settingValue(rows, DeployTimeoutKey))
	if err != nil {
		return nil, Usagef("%s: %v", DeployTimeoutKey, err).WithNext("rhiza daemon settings")
	}
	interval := deployInterval
	if ctx.Getenv != nil {
		if d, err := time.ParseDuration(ctx.Getenv("RHIZA_DEPLOY_INTERVAL")); err == nil && d > 0 {
			interval = d
		}
	}
	resp := ShipDeployResponse{URL: url, Timeout: timeout.String()}
	deadline := time.Now().Add(timeout)
	for {
		resp.Polls++
		resp.Last = deployProbe(url)
		if resp.Last == "200 OK" {
			resp.Healthy = true
			return resp, nil
		}
		if ctx.Stderr != nil {
			fmt.Fprintf(ctx.Stderr, "deploy check %d: %s\n", resp.Polls, resp.Last)
		}
		if !time.Now().Add(interval).Before(deadline) {
			return resp, nil
		}
		time.Sleep(interval)
	}
}

// deployProbe asks url once and says what came back: the HTTP status line
// or the error.
func deployProbe(url string) string {
	c, cancel := context.WithTimeout(context.Background(), deployReqTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(c, http.MethodGet, url, nil)
	if err != nil {
		return err.Error()
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err.Error()
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		return "200 OK"
	}
	return fmt.Sprintf("HTTP %d", res.StatusCode)
}
