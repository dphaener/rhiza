# step-rerun-classified-verified

evidence · `01M77ESNF44RN64HZVJDQJB079` · actor [claude-code-305277cd](../actor/claude-code-305277cd.md)

## result

PASS

## statement

Verified at fb1045461. Diff: stepRunner.Command (internal/api/run.go) now reruns a failed command step once only when Ctx.Classify (new optional func in api.go) returns transient; nil classifier means no rerun, StopExit is never rerun, and the second result stands. vanishedGoWork and its regexp are removed, and the old capability step-retry-vanished-go-work no longer resolves in the to-be view. Ran go test -trimpath ./internal/api -run TestCommandRerunsOnTransient -v -count=1: PASS for 'transient then passes' (2 runs, exit 0, 1 rerun line), 'transient twice stands' (2 runs, exit 1), 'real not rerun' (1 run, exit 1, no rerun line) and 'no classifier not rerun' (1 run). go vet ./internal/api and gofmt -l internal/api were clean. Scoped -run 'Step|TestRun' tests all pass, but the package's root-state guard flags .rhiza-root/usage.live.json as changed; that also happens without the new test and matches the known fake-usage leak, so it is not from this leaf. A full internal/api run hit the 600s timeout and later an outside SIGTERM under VM load ~7, with no FAIL lines.

## Links

### Out

- verifies [step-rerun-classified](../capability/rhiza/workflows/workflow-runs/step-rerun-classified.md)

## History

- `01M763Z4X9Q90ASH2K1BF7PZPS` Failures classified by a configurable LLM: a machine-wide setting names the classifier command (headless Claude with Haiku by default, any LLM CLI allowed); the watchdog and the step runner retry a failed step only when it says transient; the hard-coded failure signatures go; ci-triage hands the fix agent the failed CI logs instead of parsing Go, RSpec, Minitest and pytest output.
