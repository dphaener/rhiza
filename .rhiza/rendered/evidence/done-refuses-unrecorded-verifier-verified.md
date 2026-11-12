# done-refuses-unrecorded-verifier-verified

evidence · `01M78520K6ECTKCWH4E2WME5NP` · actor [claude-code-e0ecc376](../actor/claude-code-e0ecc376.md)

## result

PASS

## statement

At e62f6b443, read the diff (internal/api/session_done.go: isLeafVerifier + verifierRecordMissing, a Conflict with next `rhiza run findings --leaf <leaf>` whose message names both `rhiza run findings` and the Evidence recipe; new internal/api/session_done_verifier_test.go). Ran `go test -trimpath ./internal/api -run SessionDone -count=1 -v`: TestSessionDoneVerifierNeedsRecord PASS. It covers each acceptance point: builder session done succeeds; verifier with nothing recorded gets a conflict naming rhiza run findings and no done mark; Evidence by another actor (bob) alone is still refused with no mark; after findings journaled on the attempt it succeeds; with Evidence by the verifier's own actor (attempt 2) it succeeds. The other SessionDone tests pass too, as do `-run 'Leaves|Verifier|Evidence'`; go vet and gofmt are clean. Scripts that run verifier sessions pass (bugfix-workflow-autonomous, run-ship-ci-fix, run-ship-ci-lint-fix, run-ship-conflict-fix, run-ship-model-conflict). plan-workflow-pr fails on the plan-reviewer prompt text at line 77. That comes from sibling leaf done-own-command-guide (c44cd76a7 changed plan.toml prompts), not from this leaf's diff.

## Links

### Out

- verifies [done-refuses-unrecorded-verifier](../capability/rhiza/workflows/workflow-runs/agent-step-done-signal/done-refuses-unrecorded-verifier.md)

## History

- `01M76XNMMENQCK9PZMGC9YBG99` A verifier, review check or critic in a run cannot say done until it has recorded its result for this visit (Evidence or findings, a findings round, verdicts), and every role is told to run session done as its own command only after the recording command succeeded. Evidence: 23 of 491 build verifier attempts ended with neither Evidence nor findings, each costing a full builder and verifier retry (perf report 2026-11-12, prompt-gap rec 2).
