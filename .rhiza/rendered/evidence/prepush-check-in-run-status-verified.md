# prepush-check-in-run-status-verified

evidence · `01M773DSX6EXA0M14STV82GNVH` · actor [claude-code-89e47f72](../actor/claude-code-89e47f72.md)

## result

PASS

## statement

Verifier claude-code-89e47f72 at aab52edf3: go test -trimpath ./internal/run/... ./internal/api -run PrePush -count=1 passed (TestPrePushFold; TestPrePushJournal runs shipPrePushCheck twice on a run whose worktree is the repo and reads back a passed entry with head, command 'rhiza check' and no output, then a failed entry with the command naming the prePushCheck and output 'FAIL generated file', LatestPrePush = failed; all existing TestShipPrePushCheck* pass). One of three api runs tripped rootguard on .rhiza-root/usage.live.json, the live Claude Code statusline file the diff does not touch; reruns were clean. go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-status|describe-golden)$' passed: rhiza run status prints '## Pre-push check' with '- failed on 2222222 (rhiza check; make self-check)' and the output block, not the older passed head, and --json carries pre_push {outcome, head, command, output}. Read the diff: TypePrePush entry folded into State.PrePushChecks, journalPrePush appends to RHIZA_RUN's run or the unfinished run with that worktree, best effort, skipped on dry run; RunInfo.PrePush documented in cli.md and describe-golden.

## Links

### Out

- verifies [prepush-check-in-run-status](../capability/rhiza/orchestration/shipping/wait-green-pr/prepush-check-in-run-status.md)

## History

- `01M76RH0YQYW1SMCEE5189QDJT` Every run checks the model and generated files locally before it pushes, whatever CI's state, so model and generated-file errors never reach CI. Owner, 2026-11-11, after 'I'm still seeing more ci failures than i would like' and the orchestrator's breakdown, chose option A: "A" (require a pre-push self-check in every run). Evidence: of ~44 PR CI runs on 2026-11-11, three failed on errors a one-minute local check catches: PR #442 lint, stale docs/reference (builders added settings without rhiza docs); agent docs plan PR #456, 'duplicate handle milestone-32'; release plan PR #457, a log entry overwritten by a shell command ('not canonical'). Today ship.gate (internal/api/ship_gate.go) runs the project's rhiza.ship.localGate (this repo: 'make check' in git config) only while CI is unstable, and the plan workflow has no gate at all. Wanted: (1) a fast pre-push check that always runs before any run pushes (plan, plan-milestone, ship-change and its wrappers such as bugfix, every push including fixer and ci-fix pushes and merges of main): rhiza check on the run worktree (model, log canonicity, stale render) plus the project's generated-file command; (2) the generated-file command is a project setting, language neutral (this repo: make self-check, which also rebuilds rhiza and checks docs/reference), so other projects set their own or none; (3) the heavier rhiza.ship.localGate keeps its current meaning (only while CI is unstable); (4) a failed pre-push check goes back to the step's fixer with the failing lines instead of pushing, within the step's cap; (5) the run journal and rhiza run status show the pre-push check's result. Coordinate with milestone 32's language-neutral settings (01M763YX).
