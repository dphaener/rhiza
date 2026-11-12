package api

import (
	"strings"
	"testing"
)

func ciHealth(t *testing.T, gate, rec string) HealthPart {
	t.Helper()
	f, _, _ := shipFixture(t)
	if gate != "" {
		f.git.out["config --get "+LocalGateKey] = gate + "\n"
	}
	if err := writeAtomic(ciStabilityPath(f.common), []byte(rec)); err != nil {
		t.Fatal(err)
	}
	return healthCI(f.ctx)
}

func TestHealthCI(t *testing.T) {
	t.Run("one infra failure and a single flake stay ok", func(t *testing.T) {
		p := ciHealth(t, "", `{"runs":[
{"head":"a","result":"green","time":"x","infra":["scripts (4/4) / UNKNOWN STEP"]},
{"head":"b","result":"green","time":"x","flakes":["TestX"]}]}`)
		if p.State != HealthOK || !strings.HasPrefix(p.Summary, "stable") {
			t.Fatalf("got %+v", p)
		}
		if len(p.Detail) != 2 || !strings.Contains(p.Detail[0], "flaky TestX: 1 flake ") ||
			!strings.Contains(p.Detail[1], "infra scripts (4/4) / UNKNOWN STEP: 1 failure ") {
			t.Errorf("detail %q", p.Detail)
		}
	})
	t.Run("a test flaking on two heads warns", func(t *testing.T) {
		p := ciHealth(t, "make check", `{"runs":[
{"head":"a","result":"green","time":"x","flakes":["TestX"]},
{"head":"b","result":"green","time":"x","flakes":["TestX"]}]}`)
		if p.State != HealthWarn || !strings.Contains(p.Summary, "TestX flaked 2 times on 2 heads") || !strings.Contains(p.Summary, "stops applying after") {
			t.Errorf("got %+v", p)
		}
	})
}

func TestHealthCIStability(t *testing.T) {
	t.Run("clean record is ok", func(t *testing.T) {
		p := ciHealth(t, "", `{"runs":[{"head":"a","result":"green","time":"x"}]}`)
		if p.State != HealthOK || p.Part != HealthCI || p.Summary != "stable" || p.Next != "" || len(p.Detail) != 0 {
			t.Errorf("got %+v", p)
		}
	})
	t.Run("flakes warn with each test and its count", func(t *testing.T) {
		p := ciHealth(t, "make check", `{"runs":[
{"head":"a","result":"green","time":"x","flakes":["TestPond","TestGarden"]},
{"head":"b","result":"green","time":"x","flakes":["TestPond"]}]}`)
		if p.State != HealthWarn || !strings.HasPrefix(p.Summary, "unstable: TestPond flaked 2 times") {
			t.Fatalf("got %+v", p)
		}
		if len(p.Detail) != 2 || !strings.Contains(p.Detail[0], "TestGarden: 1 flake ") || !strings.Contains(p.Detail[1], "TestPond: 2 flakes ") {
			t.Errorf("detail %q", p.Detail)
		}
		if strings.Contains(p.Summary, "no local gate") || p.Next != "" {
			t.Errorf("gate warning with a gate set: %+v", p)
		}
	})
	t.Run("infra failures are stable and not listed as flakes", func(t *testing.T) {
		p := ciHealth(t, "", `{"runs":[
{"head":"a","result":"green","time":"x","flakes":["ci / UNKNOWN STEP"]},
{"head":"b","result":"green","time":"x","flakes":["ci / UNKNOWN STEP","CI run 123"]}]}`)
		if p.State != HealthOK || !strings.HasPrefix(p.Summary, "stable: ignored infra failures") || p.Next != "" || strings.Contains(p.Summary, "gate") {
			t.Fatalf("got %+v", p)
		}
		if len(p.Detail) != 2 || !strings.Contains(p.Detail[0], "infra CI run 123: 1 failure ") ||
			!strings.Contains(p.Detail[1], "infra ci / UNKNOWN STEP: 2 failures ") ||
			!strings.Contains(p.Detail[1], "do not make CI unstable") {
			t.Errorf("detail %q", p.Detail)
		}
	})
	t.Run("main red warns", func(t *testing.T) {
		p := ciHealth(t, "make check", `{"runs":[{"head":"m1","main":true,"result":"red","time":"x"}]}`)
		if p.State != HealthWarn || !strings.Contains(p.Summary, "main's latest CI result is red") {
			t.Errorf("got %+v", p)
		}
	})
	t.Run("unstable with no gate warns of it", func(t *testing.T) {
		p := ciHealth(t, "", `{"runs":[{"head":"m1","main":true,"result":"red","time":"x"}]}`)
		if p.State != HealthWarn || !strings.Contains(p.Summary, "no local gate is set") || !strings.Contains(p.Next, LocalGateKey) {
			t.Errorf("got %+v", p)
		}
	})
	t.Run("stable with no gate says nothing of it", func(t *testing.T) {
		p := ciHealth(t, "", `{"runs":[]}`)
		if strings.Contains(p.Summary, "gate") || p.Next != "" {
			t.Errorf("got %+v", p)
		}
	})
}
