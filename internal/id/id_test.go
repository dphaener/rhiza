package id

import (
	"bytes"
	"errors"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

const valid = "01ARZ3NDEKTSV4RRFFQ69G5FAV"

func TestParseValid(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{"typical", valid},
		{"zero", "00000000000000000000000000"},
		{"max", "7ZZZZZZZZZZZZZZZZZZZZZZZZZ"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.in)
			if err != nil {
				t.Fatalf("Parse(%q): %v", tt.in, err)
			}
			if s := got.String(); s != tt.in {
				t.Fatalf("round trip: got %q, want %q", s, tt.in)
			}
		})
	}
}

func TestParseInvalid(t *testing.T) {
	tests := []struct {
		name string
		in   string
	}{
		{"empty", ""},
		{"lowercase", strings.ToLower(valid)},
		{"mixed_case", "01ARZ3NDEKTSV4RRFFQ69G5FAv"},
		{"wrong_length_25", valid[:25]},
		{"wrong_length_27", valid + "0"},
		{"bad_char_I", "01ARZ3NDEKTSV4RRFFQ69G5FAI"},
		{"bad_char_L", "01ARZ3NDEKTSV4RRFFQ69G5FAL"},
		{"bad_char_O", "01ARZ3NDEKTSV4RRFFQ69G5FAO"},
		{"bad_char_U", "01ARZ3NDEKTSV4RRFFQ69G5FAU"},
		{"bad_char_punct", "01ARZ3NDEKTSV4RRFFQ69G5FA-"},
		{"overflow", "8ZZZZZZZZZZZZZZZZZZZZZZZZZ"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.in)
			if !errors.Is(err, ErrInvalid) {
				t.Fatalf("Parse(%q) = %v, %v; want ErrInvalid", tt.in, got, err)
			}
			if !got.IsZero() {
				t.Fatalf("Parse(%q) returned non-zero ID on error", tt.in)
			}
		})
	}
}

func TestTextRoundTrip(t *testing.T) {
	want, err := Parse(valid)
	if err != nil {
		t.Fatal(err)
	}
	b, err := want.MarshalText()
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != valid {
		t.Fatalf("MarshalText = %q, want %q", b, valid)
	}
	var got ID
	if err := got.UnmarshalText(b); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("UnmarshalText = %v, want %v", got, want)
	}
	if err := got.UnmarshalText([]byte(strings.ToLower(valid))); !errors.Is(err, ErrInvalid) {
		t.Fatalf("UnmarshalText(lowercase) = %v, want ErrInvalid", err)
	}
}

func TestCompareMatchesStringOrder(t *testing.T) {
	g := NewGenerator(time.Now, rand.New(rand.NewSource(1)))
	ids := make([]ID, 200)
	for i := range ids {
		if i%3 == 0 {
			// Vary the time too, so ordering spans timestamps and entropy.
			g = NewGenerator(func() time.Time { return time.UnixMilli(int64(rand.Intn(1 << 40))) }, rand.New(rand.NewSource(int64(i))))
		}
		id, err := g.New()
		if err != nil {
			t.Fatal(err)
		}
		ids[i] = id
	}
	for _, a := range ids {
		for _, b := range ids {
			want := strings.Compare(a.String(), b.String())
			if got := Compare(a, b); got != want {
				t.Fatalf("Compare(%v, %v) = %d, want %d", a, b, got, want)
			}
		}
	}
}

func TestIsZero(t *testing.T) {
	var z ID
	if !z.IsZero() {
		t.Fatal("zero value: IsZero = false")
	}
	id, err := Default().New()
	if err != nil {
		t.Fatal(err)
	}
	if id.IsZero() {
		t.Fatal("generated ID: IsZero = true")
	}
}

func TestGeneratorMonotonicSameMillisecond(t *testing.T) {
	fixed := time.UnixMilli(1_700_000_000_000)
	g := NewGenerator(func() time.Time { return fixed }, rand.New(rand.NewSource(42)))
	var prev ID
	for i := 0; i < 1000; i++ {
		id, err := g.New()
		if err != nil {
			t.Fatalf("call %d: %v", i, err)
		}
		if i > 0 && Compare(prev, id) >= 0 {
			t.Fatalf("call %d: %v not after %v", i, id, prev)
		}
		prev = id
	}
}

func TestGeneratorDeterministic(t *testing.T) {
	fixed := time.UnixMilli(1_700_000_000_000)
	gen := func() []ID {
		g := NewGenerator(func() time.Time { return fixed }, rand.New(rand.NewSource(7)))
		out := make([]ID, 10)
		for i := range out {
			id, err := g.New()
			if err != nil {
				t.Fatal(err)
			}
			out[i] = id
		}
		return out
	}
	a, b := gen(), gen()
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("index %d: %v != %v", i, a[i], b[i])
		}
	}
}

func TestGeneratorConcurrent(t *testing.T) {
	const workers, each = 8, 500
	g := Default()
	var (
		mu  sync.Mutex
		all []ID
		wg  sync.WaitGroup
	)
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			local := make([]ID, 0, each)
			for i := 0; i < each; i++ {
				id, err := g.New()
				if err != nil {
					t.Error(err)
					return
				}
				local = append(local, id)
			}
			mu.Lock()
			all = append(all, local...)
			mu.Unlock()
		}()
	}
	wg.Wait()
	if len(all) != workers*each {
		t.Fatalf("got %d IDs, want %d", len(all), workers*each)
	}
	sort.Slice(all, func(i, j int) bool { return Compare(all[i], all[j]) < 0 })
	for i := 1; i < len(all); i++ {
		if all[i] == all[i-1] {
			t.Fatalf("duplicate ID %v", all[i])
		}
	}
}

func TestGeneratorEntropyFailure(t *testing.T) {
	g := NewGenerator(time.Now, bytes.NewReader(nil))
	if _, err := g.New(); err == nil {
		t.Fatal("New with exhausted entropy: want error")
	}
}
