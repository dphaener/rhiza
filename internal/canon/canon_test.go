package canon

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

var goldenCases = []struct {
	name string
	in   string
	want string
}{
	{"null", `null`, "null\n"},
	{"true", `true`, "true\n"},
	{"false", ` false `, "false\n"},
	{"zero", `0`, "0\n"},
	{"negative zero", `-0`, "0\n"},
	{"negative", `-42`, "-42\n"},
	{"int64 max", `9223372036854775807`, "9223372036854775807\n"},
	{"int64 min", `-9223372036854775808`, "-9223372036854775808\n"},
	{"surrounding whitespace", " \t\r\n{\"a\":1}\n\n ", "{\n  \"a\": 1\n}\n"},
	{"empty object", `{}`, "{}\n"},
	{"empty object with space", `{ }`, "{}\n"},
	{"empty array", `[ ]`, "[]\n"},
	{"sorted keys", `{"b":1,"a":2,"c":3}`, "{\n  \"a\": 2,\n  \"b\": 1,\n  \"c\": 3\n}\n"},
	{"byte-wise key order", `{"é":1,"z":2,"Z":3,"a":4,"😀":5}`,
		"{\n  \"Z\": 3,\n  \"a\": 4,\n  \"z\": 2,\n  \"é\": 1,\n  \"😀\": 5\n}\n"},
	{"key order uses unescaped bytes", `{"\u0062":1,"a":2}`, "{\n  \"a\": 2,\n  \"b\": 1\n}\n"},
	{"nesting", `{"b":[1,{"a":true}],"a":{},"c":[]}`,
		"{\n  \"a\": {},\n  \"b\": [\n    1,\n    {\n      \"a\": true\n    }\n  ],\n  \"c\": []\n}\n"},
	{"array of arrays", `[[],[null]]`, "[\n  [],\n  [\n    null\n  ]\n]\n"},
	{"short escapes", `"\"\\\b\f\n\r\t"`, "\"\\\"\\\\\\b\\f\\n\\r\\t\"\n"},
	{"other controls lower hex", `"\u0000\u0001\u001F\u000b"`, "\"\\u0000\\u0001\\u001f\\u000b\"\n"},
	{"solidus raw", `"a\/b/c"`, "\"a/b/c\"\n"},
	{"html raw", `"<>&\u003c\u003E\u0026"`, "\"<>&<>&\"\n"},
	{"del raw", "\"\x7f\\u007f\"", "\"\x7f\x7f\"\n"},
	{"e acute raw", `"é\u00e9"`, "\"éé\"\n"},
	{"emoji raw", `"😀\ud83d\ude00"`, "\"😀😀\"\n"},
	{"u2028 u2029 raw", `"\u2028\u2029` + "\u2028\"", "\"\u2028\u2029\u2028\"\n"},
	{"no normalisation", "\"e\u0301\"", "\"e\u0301\"\n"},
	{"control escape in key", `{"a\nb":1}`, "{\n  \"a\\nb\": 1\n}\n"},
}

var invalidCases = []struct {
	name string
	in   string
}{
	{"empty", ``},
	{"only whitespace", "  \n"},
	{"duplicate key", `{"a":1,"a":2}`},
	{"duplicate key after unescape", `{"a":1,"\u0061":2}`},
	{"fraction", `1.5`},
	{"exponent", `1e3`},
	{"zero fraction", `1.0`},
	{"capital exponent", `1E3`},
	{"above int64", `9223372036854775808`},
	{"below int64", `-9223372036854775809`},
	{"leading zero", `01`},
	{"plus sign", `+1`},
	{"bare minus", `-`},
	{"invalid utf8 byte", "\"\xff\""},
	{"invalid utf8 in key", "{\"\xc3\":1}"},
	{"lone high surrogate", `"\ud800"`},
	{"lone low surrogate", `"\udc00"`},
	{"high then non-low", `"\ud800\u0041"`},
	{"trailing garbage", `{} x`},
	{"two values", `1 2`},
	{"raw control in string", "\"a\nb\""},
	{"bad escape", `"\x"`},
	{"short unicode escape", `"\u12"`},
	{"unterminated string", `"abc`},
	{"unterminated object", `{"a":1`},
	{"trailing comma", `[1,]`},
	{"missing colon", `{"a" 1}`},
	{"non-string key", `{1:2}`},
	{"bad literal", `tru`},
	{"too deep", strings.Repeat("[", maxDepth+1) + strings.Repeat("]", maxDepth+1)},
}

func TestCanonicalizeGolden(t *testing.T) {
	for _, tc := range goldenCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Canonicalize([]byte(tc.in))
			if err != nil {
				t.Fatalf("Canonicalize(%q): %v", tc.in, err)
			}
			if string(got) != tc.want {
				t.Fatalf("Canonicalize(%q)\n got: %q\nwant: %q", tc.in, got, tc.want)
			}
			if err := Check(got); err != nil {
				t.Fatalf("Check(output): %v", err)
			}
		})
	}
}

func TestCanonicalizeInvalid(t *testing.T) {
	for _, tc := range invalidCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Canonicalize([]byte(tc.in))
			if !errors.Is(err, ErrInvalid) {
				t.Fatalf("Canonicalize(%q) = %q, %v; want ErrInvalid", tc.in, got, err)
			}
			if err := Check([]byte(tc.in)); !errors.Is(err, ErrInvalid) {
				t.Fatalf("Check(%q) = %v; want ErrInvalid", tc.in, err)
			}
		})
	}
}

func TestMaxDepthAccepted(t *testing.T) {
	in := strings.Repeat("[", maxDepth) + strings.Repeat("]", maxDepth)
	if _, err := Canonicalize([]byte(in)); err != nil {
		t.Fatalf("depth %d: %v", maxDepth, err)
	}
}

func TestCheck(t *testing.T) {
	accept := []string{
		"{\n  \"a\": 1,\n  \"b\": [\n    true\n  ]\n}\n",
		"{}\n",
		"\"é<>&\"\n",
		"-9223372036854775808\n",
	}
	for _, in := range accept {
		if err := Check([]byte(in)); err != nil {
			t.Errorf("Check(%q) = %v; want nil", in, err)
		}
	}

	reject := []struct{ name, in string }{
		{"missing final newline", "{\n  \"a\": 1\n}"},
		{"two final newlines", "{\n  \"a\": 1\n}\n\n"},
		{"tab indent", "{\n\t\"a\": 1\n}\n"},
		{"unsorted keys", "{\n  \"b\": 1,\n  \"a\": 2\n}\n"},
		{"escaped e acute", "\"\\u00e9\"\n"},
		{"compact", "{\"a\":1}\n"},
		{"negative zero", "-0\n"},
		{"upper hex escape", "\"\\u001F\"\n"},
		{"escaped solidus", "\"\\/\"\n"},
		{"leading whitespace", " 1\n"},
	}
	for _, tc := range reject {
		t.Run(tc.name, func(t *testing.T) {
			err := Check([]byte(tc.in))
			if !errors.Is(err, ErrNotCanonical) {
				t.Fatalf("Check(%q) = %v; want ErrNotCanonical", tc.in, err)
			}
			if errors.Is(err, ErrInvalid) {
				t.Fatalf("Check(%q) = %v; must not also wrap ErrInvalid", tc.in, err)
			}
		})
	}
}

func TestMarshal(t *testing.T) {
	v := map[string]any{
		"z":    []string{"<b>", "a&b"},
		"a":    int64(-9223372036854775808),
		"m":    map[string]int{},
		"list": []int{},
	}
	got, err := Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	want := "{\n  \"a\": -9223372036854775808,\n  \"list\": [],\n  \"m\": {},\n  \"z\": [\n    \"<b>\",\n    \"a&b\"\n  ]\n}\n"
	if string(got) != want {
		t.Fatalf("Marshal\n got: %q\nwant: %q", got, want)
	}

	if _, err := Marshal(1.5); !errors.Is(err, ErrInvalid) {
		t.Fatalf("Marshal(1.5) err = %v; want ErrInvalid", err)
	}
	if _, err := Marshal(func() {}); err == nil {
		t.Fatal("Marshal(func) succeeded; want error")
	}
}

func TestDeterministic(t *testing.T) {
	in := []byte(`{"b":{"y":[1,2],"x":"é"},"a":null}`)
	first, err := Canonicalize(in)
	if err != nil {
		t.Fatal(err)
	}
	for range 50 {
		got, err := Canonicalize(in)
		if err != nil || !bytes.Equal(got, first) {
			t.Fatalf("Canonicalize not deterministic: %q, %v", got, err)
		}
	}
}

func FuzzCanonicalize(f *testing.F) {
	for _, tc := range goldenCases {
		f.Add([]byte(tc.in))
		f.Add([]byte(tc.want))
	}
	for _, tc := range invalidCases {
		f.Add([]byte(tc.in))
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		out, err := Canonicalize(data)
		checkErr := Check(data)
		if err != nil {
			if !errors.Is(err, ErrInvalid) {
				t.Fatalf("Canonicalize error does not wrap ErrInvalid: %v", err)
			}
			if !errors.Is(checkErr, ErrInvalid) {
				t.Fatalf("Canonicalize rejects but Check = %v", checkErr)
			}
			return
		}
		if (checkErr == nil) != bytes.Equal(out, data) {
			t.Fatalf("Check(%q) = %v disagrees with Canonicalize = %q", data, checkErr, out)
		}
		if err := Check(out); err != nil {
			t.Fatalf("Check(Canonicalize(%q)) = %v", data, err)
		}
		again, err := Canonicalize(out)
		if err != nil || !bytes.Equal(again, out) {
			t.Fatalf("not idempotent: %q -> %q -> %q, %v", data, out, again, err)
		}
	})
}
