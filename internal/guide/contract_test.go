package guide

import (
	"fmt"
	"strings"
	"testing"
)

// roleTopics are the topics an agent loads for its role in a run.
var roleTopics = []string{
	"advisor", "arbiter", "builder", "verifier",
	"check-compliance", "check-model", "check-qa", "check-quality",
	"critic", "fixer", "ci-fix", "gate-fix", "conflict-fix",
	"plan-critic", "plan-reviewer", "plan-reviser", "planner",
}

// missingContractRef returns an error naming each role topic of texts that
// does not name rhiza guide contract.
func missingContractRef(texts []Text) error {
	var bad []string
	for _, t := range texts {
		for _, r := range roleTopics {
			if t.Name == r && !strings.Contains(t.Body, "rhiza guide contract") {
				bad = append(bad, r)
			}
		}
	}
	if len(bad) > 0 {
		return fmt.Errorf("role topics without rhiza guide contract: %s", strings.Join(bad, ", "))
	}
	return nil
}

func TestContract(t *testing.T) {
	c, ok := Lookup("contract")
	if !ok {
		t.Fatal("no contract topic")
	}
	for _, want := range []string{
		"`Change:`", "`Actor:`", "`Leaf:`", "Evidence", "`rhiza run findings`", "`rhiza session done`",
		"gate",
		"{brief}", "{slice}", "{handback}", "{change}", "{run}", "{leaf}", "{attempt}", "{worktree}",
		"RHIZA_BRIEF", "RHIZA_SLICE", "RHIZA_HANDBACK", "RHIZA_CHANGE", "RHIZA_RUN", "RHIZA_LEAF", "RHIZA_ATTEMPT",
		"brief.md", "slice.json", "handback.md",
	} {
		if !strings.Contains(c.Body, want) {
			t.Errorf("contract lacks %s", want)
		}
	}
	if err := checkCap(c); err != nil {
		t.Error(err)
	}
	if err := missingContractRef(List()); err != nil {
		t.Error(err)
	}
}

func TestContractRefFailureNamesTopic(t *testing.T) {
	var texts []Text
	for _, tx := range List() {
		if tx.Name == "verifier" {
			tx.Body = strings.ReplaceAll(tx.Body, "rhiza guide contract", "")
		}
		texts = append(texts, tx)
	}
	err := missingContractRef(texts)
	if err == nil || !strings.Contains(err.Error(), "verifier") {
		t.Fatalf("missingContractRef = %v, want an error naming verifier", err)
	}
}
