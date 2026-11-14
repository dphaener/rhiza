# code-check-runs-lints-qa-1

evidence · `01M7B95FGHJ8R4V7AV9G6HTDYY` · actor [claude-code-59e3a197](../actor/claude-code-59e3a197.md)

## result

PASS

## statement

Built ./cmd/rhiza to a scratch path and ran it in a hermetic scratch repo with a run stopped at a wait step: leaf acceptance naming internal/nope.go, an unlinked Decision and a branch deleting the only file of a landed capability. With alpha lacking verifier Evidence, run code-check exited 1 with only the Evidence cause and listed 6 lint hits (path, test, script, unapplied-decision, removed-code-capability) in human output and --json; the journal entry carried them. After Evidence it exited 0 with passed true and the same hits; resuming the run, run status (human and --json) showed the hits under code_checks for every entry. --json and api call run.code-check gave identical bytes; dry-run lints were byte-stable across runs. Scoped TestScripts/(review-code-check|describe-golden) and internal/run, internal/api -run CodeCheck pass.

## Links

### Out

- verifies [code-check-runs-lints](../capability/rhiza/workflows/built-in-workflows/built-in-review/code-check-runs-lints.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
