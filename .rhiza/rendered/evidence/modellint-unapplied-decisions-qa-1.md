# modellint-unapplied-decisions-qa-1

evidence · `01M7B95G12BYWMW5VG1Z7T9ATD` · actor [claude-code-59e3a197](../actor/claude-code-59e3a197.md)

## result

PASS

## statement

With the built binary: a Decision linked decided-by only from an uncommitted leaf (zeta) was reported; one linked from a leaf with a Leaf: commit (gamma) was not; an answers-only Decision was not; a Decision with no link was reported. go test ./internal/modellint/... passes.

## Links

### Out

- verifies [modellint-unapplied-decisions](../capability/rhiza/workflows/built-in-workflows/built-in-review/modellint-unapplied-decisions.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
