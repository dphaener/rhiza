# workflow-new-placeholder-verified

evidence · `01M77ZCV6F5YT5D38M49ZDN274` · actor [claude-code-115c2567](../actor/claude-code-115c2567.md)

## result

PASS

## statement

Verifier at e8e7b194c: read the diff (workflow_new.go skeleton with RHIZA-PLACEHOLDER test command that exits 1, commented fields, Go/Ruby/JS/TS/Python/Rust examples, commented [[cleanup]] and [env]; workflow_check.go placeholderWarnings appends a warning and never touches Problems/Valid). Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(workflow-new|change-workflow-pin|workflow-check)$' with GOTMPDIR/TMPDIR in a neutral temp dir: all PASS. workflow-new.txtar shows workflow new skel exit 0, no 'make test', all language examples and commented cleanup/env present, workflow check skel exit 0 printing 'warning: step test still runs the placeholder', --json valid true, a fork from ship-change with no placeholder warning, and after replacing run with go test ./... check says valid with no warning. go test ./internal/workflow PASS; go test ./internal/api -run Workflow PASS. Note: the script's ! stdout 'placeholder' on 'workflow check mine' matches the printed $WORK path if TMPDIR contains 'placeholder' (as this run's scratchpad did); it is environment-only, the repo's own temp paths do not.

## Links

### Out

- verifies [workflow-new-placeholder](../capability/rhiza/workflows/workflow-authoring/workflow-new-placeholder.md)

## History

- `01M763Z8TGV8GJK2E621JS6HYZ` Language-neutral workflow scaffold, docs and examples: rhiza workflow new writes a placeholder skeleton with extensive comments and example commands for popular languages; the docs explain the cleanup phase and workflow environment and carry example workflows for several languages.
