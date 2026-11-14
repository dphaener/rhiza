package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const apiTestToken = "s3cret-token-value-1234"

// daemon.set accepts rhiza.web.token; the daemon's web handler gates POSTs on
// the token it reads from that setting; no read path returns its value.
func TestWebTokenSetting(t *testing.T) {
	ctx, _ := daemonRepo(t)
	h := WebProjectHandler(ctx)
	post := func(hdr string) int {
		r := httptest.NewRequest(http.MethodPost, "/api/daemon.settings", nil)
		if hdr != "" {
			r.Header.Set("Authorization", hdr)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		return rec.Code
	}
	if c := post(""); c != http.StatusForbidden {
		t.Fatalf("POST with no token set = %d, want 403", c)
	}
	if _, err := runDaemonSet(ctx, DaemonSetRequest{Key: "rhiza.web.token", Value: "short"}); err == nil {
		t.Fatal("daemon set accepted a short token")
	}
	resp, err := runDaemonSet(ctx, DaemonSetRequest{Key: "rhiza.web.token", Value: apiTestToken})
	if err != nil {
		t.Fatal(err)
	}
	if r := resp.(DaemonSetResponse); r.Value != "set" || strings.Contains(r.Value, apiTestToken) {
		t.Errorf("daemon set replied %+v", r)
	}
	// The gate reads the setting at each request.
	if c := post("Bearer wrong"); c != http.StatusUnauthorized {
		t.Errorf("POST with a wrong token = %d, want 401", c)
	}
	if c := post("Bearer " + apiTestToken); c == http.StatusUnauthorized || c == http.StatusForbidden {
		t.Errorf("POST with the token = %d, want it through the gate", c)
	}
	// No read path shows the value.
	out, err := callOp(t, ctx, "daemon.settings", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	rows := out.(DaemonSettingsResponse).Settings
	if v := settingValue(rows, "rhiza.web.token"); v != "set" {
		t.Errorf("daemon.settings token value = %q, want set", v)
	}
	for _, s := range (daemonConf{rows: rows}).stateRows() {
		if strings.Contains(s.Value, apiTestToken) {
			t.Errorf("daemon state row %s carries the token", s.Key)
		}
	}
	for _, target := range []string{"/api/daemon.settings", "/api/daemon.status"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, target, nil))
		if strings.Contains(rec.Body.String(), apiTestToken) {
			t.Errorf("GET %s returns the token:\n%s", target, rec.Body)
		}
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/daemon.settings", nil))
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), `"rhiza.web.token"`) {
		t.Errorf("GET /api/daemon.settings = %d %s", rec.Code, rec.Body)
	}
}
