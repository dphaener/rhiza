# modellint-branch-removed-code-qa-1

evidence · `01M7B95G8WE9HT9PV48S56CE29` · actor [claude-code-59e3a197](../actor/claude-code-59e3a197.md)

## result

PASS

## statement

With the built binary: a landed capability naming the directory olddir/, whose files the branch deleted, was reported removed-code-capability with ref olddir; a landed capability naming one deleted file and one still-present file was not reported; a landed capability naming legacy/old.txt deleted on the branch was reported. go test ./internal/modellint/... passes.

## Links

### Out

- verifies [modellint-branch-removed-code](../capability/modellint-branch-removed-code.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
