# fixer-lint-hits-guide-qa-1

evidence · `01M7B95GPGTFPZ50G6MGG95AM3` · actor [claude-code-59e3a197](../actor/claude-code-59e3a197.md)

## result

PASS

## statement

Built binary: rhiza guide fixer has the 'A lint-fix fixer' section naming rhiza run lint-check --json, the three fixes (rhiza set the text, retire the capability, link decided-by) and listing a wrongly judged hit as unfixed; the Checks: and never-push rules are still present. go test ./internal/guide/... passes.

## Links

### Out

- verifies [fixer-lint-hits-guide](../capability/fixer-lint-hits-guide.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
