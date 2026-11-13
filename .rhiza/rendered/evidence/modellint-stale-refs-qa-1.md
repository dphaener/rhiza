# modellint-stale-refs-qa-1

evidence · `01M7B95FTE7K9BQYDSGTS79JFG` · actor [claude-code-59e3a197](../actor/claude-code-59e3a197.md)

## result

PASS

## statement

With the built binary, code-check on a scratch change set linted only its own nodes: a missing path, a missing Go test and a missing TestScripts/<name> each gave one hit; 'and/or', an https URL, internal/real.go:12, ./internal/... and backticked present paths gave none; a capability whose two paths are both gone was also reported unretired-capability. Rough edge: a backticked placeholder TestXxx in this change set's statement is linted as a missing test (finding recorded, minor). go test ./internal/modellint/... passes.

## Links

### Out

- verifies [modellint-stale-refs](../capability/rhiza/workflows/built-in-workflows/built-in-review/modellint-stale-refs.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
