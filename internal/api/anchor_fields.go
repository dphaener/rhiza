package api

import (
	"encoding/json"
	"errors"
	"slices"

	"github.com/dphaener/rhiza/internal/anchor"
)

// anchorFingerprint applies the anchor rules to a write. all is the anchor's
// fields as they will stand; given is what this write sets. Where all has a
// path, it and the symbol must resolve against the working tree under root.
// The fingerprint is computed into the returned set when this write gave no
// fingerprint of its own and either it is new (add) or it set a path or
// symbol (set). drop reports that all has no path, so a stale fingerprint
// must go.
func anchorFingerprint(root string, all, given map[string]json.RawMessage, unset []string, isSet bool) (set map[string]json.RawMessage, drop bool, err error) {
	set = given
	path, symbol := textField(all, "path"), textField(all, "symbol")
	if path == "" {
		if symbol != "" {
			return nil, false, Usagef("an anchor symbol needs a path").WithNext("rhiza describe set")
		}
		_, hasFP := all["fingerprint"]
		_, givenFP := given["fingerprint"]
		return set, isSet && hasFP && !givenFP, nil
	}
	_, hasFP := given["fingerprint"]
	_, hasPath := given["path"]
	_, hasSymbol := given["symbol"]
	// Unsetting the symbol changes what the fingerprint covers.
	changed := !isSet || hasPath || hasSymbol || slices.Contains(unset, "symbol")
	// A set that changes neither needs no new resolve.
	if isSet && !changed && !hasFP {
		return set, false, nil
	}
	fp, err := anchor.Fingerprint(root, path, symbol)
	if err != nil {
		return nil, false, anchorError(err)
	}
	if hasFP || !changed {
		return set, false, nil
	}
	out := make(map[string]json.RawMessage, len(given)+1)
	for k, v := range given {
		out[k] = v
	}
	raw, _ := json.Marshal(fp)
	out["fingerprint"] = raw
	return out, false, nil
}

// textField is the string value of field name, or "" when absent or not a
// string.
func textField(fields map[string]json.RawMessage, name string) string {
	var s string
	if raw, ok := fields[name]; ok && json.Unmarshal(raw, &s) == nil {
		return s
	}
	return ""
}

// anchorError turns a resolve failure into a usage error that says what to
// fix.
func anchorError(err error) error {
	next := "rhiza describe set"
	switch {
	case errors.Is(err, anchor.ErrNoPath), errors.Is(err, anchor.ErrBadPath):
		next = "rhiza set <anchor> --field path=<path relative to the project root>"
	case errors.Is(err, anchor.ErrNoSymbol), errors.Is(err, anchor.ErrBadSymbol):
		next = "rhiza set <anchor> --field symbol=<Name or Type.Method> (or --unset symbol)"
	}
	return Usagef("anchor: %v", err).WithNext(next)
}
