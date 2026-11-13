# review-lint-fix-route-qa-1

evidence · `01M7B95GDDQC2FJA08EGEH2QZ2` · actor [claude-code-59e3a197](../actor/claude-code-59e3a197.md)

## result

PASS

## statement

Built binary: workflow describe review shows code -> lints, lints pass -> review, fail -> lint-fix, lint-fix agent max_visits 1 pass -> review; workflow check review is valid. run lint-check exited 0 with no code-check entry yet, exit 2 usage with no RHIZA_RUN, and exit 1 listing all hits (human and --json) after a code-check with hits; as a workflow command step after code it failed the step. go test ./internal/workflow/... and ./internal/api -run LintCheck pass.

## Links

### Out

- verifies [review-lint-fix-route](../capability/review-lint-fix-route.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
