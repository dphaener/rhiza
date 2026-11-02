package daemon

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// MainStaleFile is the name, under Dir, of the record of the last ship
// close-out that could not fast-forward the main checkout.
const MainStaleFile = "main-stale.json"

// MainStale is why the main checkout could not be brought up to date after
// a merge: ship close-out records it, rhiza health and the inbox show it,
// and the next close-out that fast-forwards main clears it.
type MainStale struct {
	// Since is when the first failure of this run of failures happened; a
	// failure that repeats keeps it, so the inbox item keeps its ID.
	Since time.Time `json:"since"`
	// Change is the change set whose close-out failed last.
	Change string `json:"change"`
	// Detail says why main was not updated.
	Detail string `json:"detail"`
}

// WriteMainStale records a failed pull, keeping the Since of an earlier
// record that is still there.
func WriteMainStale(common string, s MainStale) error {
	if prev, ok := ReadMainStale(common); ok && !prev.Since.IsZero() {
		s.Since = prev.Since
	}
	b, err := json.Marshal(s)
	if err != nil {
		return err
	}
	path := filepath.Join(Dir(common), MainStaleFile)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(b, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// ReadMainStale is the recorded pull failure, and false if there is none
// or it does not parse.
func ReadMainStale(common string) (MainStale, bool) {
	var s MainStale
	b, err := os.ReadFile(filepath.Join(Dir(common), MainStaleFile))
	if err != nil || json.Unmarshal(b, &s) != nil {
		return MainStale{}, false
	}
	return s, true
}

// ClearMainStale removes the record; none is not an error.
func ClearMainStale(common string) error {
	err := os.Remove(filepath.Join(Dir(common), MainStaleFile))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return nil
}
