# review-guide-verified

evidence · `01M5ZTESYCKGQYBMT8P8EEX7ZH` · actor [claude-code-0e5a9d92-verifier](../actor/claude-code-0e5a9d92-verifier.md)

## result

PASS

## statement

Re-ran at 64215d0 in a detached worktree, RHIZA_* unset. guide ls: guide.txtar asserts 12 entries incl '- review topic ... Reviewing a built change set', TestListOrder has review; built binary's guide ls shows 'review topic 650 tokens'. TestReviewText asserts the five checks (compliance on each leaf's acceptance, quality, qa with built binary/fakes alone do not count, scope inside the intent and the leaves, model), merge-base with the default branch, blocker/major/minor/nit, rhiza run findings --check "$RHIZA_CHECK" <file> with {"findings":[{severity,leaf,place,reason}]} and {"findings":[]}, rhiza run verdict --check with {"verdicts":[{finding,verdict,reason}]} confirmed or dropped, 'Never add a finding', Leaf: <leaf> and Review: <run> trailers, QA Evidence on each leaf exercised with a verifies link, missing constraints and discovered-from draft change set; guide.txtar asserts the same through the CLI. Formats match checkFindingsFile and run.Verdict in internal/api/run_review.go (DisallowUnknownFields). TestCaps: review 650 tokens <= TopicCap 800. TestGuideInstallReview and guide-install.txtar: install writes skills/rhiza-review/SKILL.md (7 files), body == rhiza guide review; hands-on install into a temp HOME wrote it and cmp matched. All PASS; go test ./... all ok.

## tests

env -u RHIZA_CHANGE -u RHIZA_ACTOR -u RHIZA_FOR -u RHIZA_RUN -u RHIZA_SESSION -u RHIZA_ROOT: go test ./internal/guide -v -count=1; go test ./internal/api -run 'GuideInstall|SkillFile' -v -count=1; go test ./cmd/rhiza -run 'TestScripts/(guide|guide-install|status)$' -v -count=1; go test ./... -count=1; go build -o $TMPDIR/rhiza-vg ./cmd/rhiza then rhiza-vg guide ls, guide review, guide install --agent claude (temp HOME); commit 64215d0; run by claude-code-0e5a9d92-verifier on 2026-10-27

## Links

### Out

- verifies [review-guide](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide.md)

## History

- `01M5WPTZ1V0JVVW8FPX0QYMRWJ` The review workflow: fresh sessions check compliance, quality, hands-on QA, scope and the model, a critic that can't add findings settles failed checks, fixes go back per leaf at most twice, and the owner signs off with rhiza run approve.
