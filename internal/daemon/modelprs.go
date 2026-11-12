package daemon

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

// ModelPRsFile is the name, under Dir, of the record of the open model-only
// pull requests the daemon could not merge.
const ModelPRsFile = "model-prs.json"

// ModelPR is one open model-only pull request that is red or conflicting:
// the daemon does not merge it, and the inbox lists it.
type ModelPR struct {
	Branch string `json:"branch"`
	Change string `json:"change"`
	// Since is when the daemon first found it stuck.
	Since  time.Time `json:"since"`
	Detail string    `json:"detail"`
}

// WriteModelPRs records the stuck model PRs, keeping the Since of one that
// was already recorded under the same branch. An empty list removes the file.
func WriteModelPRs(common string, prs []ModelPR) error {
	path := filepath.Join(Dir(common), ModelPRsFile)
	if len(prs) == 0 {
		if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		return nil
	}
	prev := map[string]time.Time{}
	for _, p := range ReadModelPRs(common) {
		prev[p.Branch] = p.Since
	}
	for i := range prs {
		if s, ok := prev[prs[i].Branch]; ok && !s.IsZero() {
			prs[i].Since = s
		}
	}
	b, err := json.Marshal(prs)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(b, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

// ReadModelPRs is the recorded stuck model PRs; none if there is no record
// or it does not parse.
func ReadModelPRs(common string) []ModelPR {
	var prs []ModelPR
	b, err := os.ReadFile(filepath.Join(Dir(common), ModelPRsFile))
	if err != nil || json.Unmarshal(b, &prs) != nil {
		return nil
	}
	return prs
}
