# prepush-check-before-push-verified

evidence · `01M772ZPQ5T6Q5N0V49P85RF1H` · actor [claude-code-4f0a24a2](../actor/claude-code-4f0a24a2.md)

## result

PASS

## statement

Verifier at 65e57e94e: ran go test -trimpath -count=1 ./internal/api -run 'TestShipPrePushCheck|TestShipLocalGate|TestWatchdogPRNoCI' -v; every test passed (TestShipPrePushCheck, ...ShipPush, ...OpenPR, ...Land, ...Rebase, ...WithLocalGate, ...Watchdog, ...HeadBuild, TestShipLocalGate, TestShipLocalGateOncePerHead, TestWatchdogPRNoCI). The package then exited FAIL only because rootguard saw ~/code/.rhiza-root/usage.live.json change. The live root daemon rewrites that file every few seconds: it also changed during a run that matched no tests. Read the diff: shipPush (open-pr, land), shipFinishRebase and watchdogMergeMain each call one helper, prePushCheck. It checks HEAD out into a temporary detached worktree, removed afterwards, and runs rhiza check there, then rhiza.ship.prePushCheck there. On failure it pushes nothing. open-pr, land and rebase return pre_push_check {head, lines} and exit 1; rebase and the watchdog reset --keep to the pre-merge head. The localGate stays gated on CI being unstable. The tests drive the ship.open-pr, ship.land and ship.rebase ops and watchdogMergeMain against each of the three failures: a committed non-canonical log entry, a committed stale render, and a failing command. They assert that origin is unchanged or that no push call was made, that the failing lines appear in the response, error and human output, and that one worktree is left. Two more cases fail with the fix only in the working tree. An untracked stray log entry passes. Passing pushes go through with CI stable and with CI unstable, and the gate runs only when CI is unstable. Also ran: go test ./cmd/rhiza -run 'TestScripts/(plan-workflow-pr|run-ship-change-rebase|ship-land-branch-build|describe-golden)$' passed, gofmt -l clean, go vet ./internal/api clean.

## Links

### Out

- verifies [prepush-check-before-push](../capability/rhiza/orchestration/shipping/wait-green-pr/prepush-check-before-push.md)

## History

- `01M76RH0YQYW1SMCEE5189QDJT` Every run checks the model and generated files locally before it pushes, whatever CI's state, so model and generated-file errors never reach CI. Owner, 2026-11-11, after 'I'm still seeing more ci failures than i would like' and the orchestrator's breakdown, chose option A: "A" (require a pre-push self-check in every run). Evidence: of ~44 PR CI runs on 2026-11-11, three failed on errors a one-minute local check catches: PR #442 lint, stale docs/reference (builders added settings without rhiza docs); agent docs plan PR #456, 'duplicate handle milestone-32'; release plan PR #457, a log entry overwritten by a shell command ('not canonical'). Today ship.gate (internal/api/ship_gate.go) runs the project's rhiza.ship.localGate (this repo: 'make check' in git config) only while CI is unstable, and the plan workflow has no gate at all. Wanted: (1) a fast pre-push check that always runs before any run pushes (plan, plan-milestone, ship-change and its wrappers such as bugfix, every push including fixer and ci-fix pushes and merges of main): rhiza check on the run worktree (model, log canonicity, stale render) plus the project's generated-file command; (2) the generated-file command is a project setting, language neutral (this repo: make self-check, which also rebuilds rhiza and checks docs/reference), so other projects set their own or none; (3) the heavier rhiza.ship.localGate keeps its current meaning (only while CI is unstable); (4) a failed pre-push check goes back to the step's fixer with the failing lines instead of pushing, within the step's cap; (5) the run journal and rhiza run status show the pre-push check's result. Coordinate with milestone 32's language-neutral settings (01M763YX).
