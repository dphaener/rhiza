package daemon

import (
	"testing"
	"time"
)

func TestParseLoadavg(t *testing.T) {
	l1, l5, err := ParseLoadavg("16.25 9.50 4.00 3/812 12345\n")
	if err != nil || l1 != 16.25 || l5 != 9.5 {
		t.Fatalf("got %v %v %v", l1, l5, err)
	}
	for _, bad := range []string{"", "1.0", "x 1.0"} {
		if _, _, err := ParseLoadavg(bad); err == nil {
			t.Errorf("ParseLoadavg(%q) did not fail", bad)
		}
	}
}

func TestParseSysctlLoadavg(t *testing.T) {
	// struct loadavg {fixpt_t ldavg[3]; long fscale}: 2.5, 1.25, 0.5 at
	// scale 2048.
	b := []byte{
		0x00, 0x14, 0, 0, 0x00, 0x0a, 0, 0, 0x00, 0x04, 0, 0, 0, 0, 0, 0,
		0x00, 0x08, 0, 0, 0, 0, 0, 0,
	}
	l1, l5, err := ParseSysctlLoadavg(b)
	if err != nil || l1 != 2.5 || l5 != 1.25 {
		t.Fatalf("got %v %v %v", l1, l5, err)
	}
	if _, _, err := ParseSysctlLoadavg(b[:12]); err == nil {
		t.Error("short vm.loadavg did not fail")
	}
	if _, _, err := ParseSysctlLoadavg(make([]byte, 24)); err == nil {
		t.Error("zero scale did not fail")
	}
}

func TestParseMemAvailable(t *testing.T) {
	in := "MemTotal:       16000000 kB\nMemFree:          100000 kB\nMemAvailable:    2097152 kB\nBuffers: 1 kB\n"
	n, err := ParseMemAvailable(in)
	if err != nil || n != 2<<30 {
		t.Fatalf("got %d %v", n, err)
	}
	if _, err := ParseMemAvailable("MemTotal: 1 kB\n"); err == nil {
		t.Error("missing MemAvailable did not fail")
	}
}

func TestEvaluate(t *testing.T) {
	now := time.Date(2026, 11, 11, 0, 0, 0, 0, time.UTC)
	r := func(l1, l5 float64, mem uint64) Reading {
		return Reading{Load1: l1, Load5: l5, LoadOK: true, MemAvail: mem, MemOK: true}
	}
	const gb = 1 << 30
	h := Evaluate(r(4, 4, 8*gb), 3, 2*gb, nil, now)
	if h == nil || h.Reason != SkipLoadHigh {
		t.Fatalf("over limit: %+v", h)
	}
	later := now.Add(time.Minute)
	h2 := Evaluate(r(1, 3.5, 8*gb), 3, 2*gb, h, later)
	if h2 == nil || !h2.Since.Equal(now) {
		t.Fatalf("hysteresis: %+v", h2)
	}
	if h3 := Evaluate(r(1, 2, 8*gb), 3, 2*gb, h2, later); h3 != nil {
		t.Fatalf("released: %+v", h3)
	}
	if h := Evaluate(r(1, 5, 8*gb), 3, 2*gb, nil, now); h != nil {
		t.Fatalf("began on 5-minute: %+v", h)
	}
	if h := Evaluate(r(1, 1, gb), 3, 2*gb, nil, now); h == nil || h.Reason != SkipMemoryLow {
		t.Fatalf("memory: %+v", h)
	}
	if h := Evaluate(r(9, 9, gb), 0, 0, nil, now); h != nil {
		t.Fatalf("disabled: %+v", h)
	}
	// A memory hold over a load hold keeps the load hold on beneath it:
	// once memory recovers, the 5-minute load still has to come down.
	m := Evaluate(r(1, 3.5, gb), 3, 2*gb, h, later)
	if m == nil || m.Reason != SkipMemoryLow || !m.Load {
		t.Fatalf("memory over load: %+v", m)
	}
	if h := Evaluate(r(1, 3.5, 8*gb), 3, 2*gb, m, later); h == nil || h.Reason != SkipLoadHigh {
		t.Fatalf("load hold dropped after memory: %+v", h)
	}
	if m := Evaluate(r(1, 3.5, gb), 3, 2*gb, nil, now); m == nil || m.Load {
		t.Fatalf("5-minute load began a load hold under memory: %+v", m)
	}
	if h := Evaluate(Reading{}, 3, 2*gb, nil, now); h != nil {
		t.Fatalf("unknown reading held: %+v", h)
	}
}
