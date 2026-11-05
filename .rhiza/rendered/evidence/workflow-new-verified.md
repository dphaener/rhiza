# workflow-new-verified

evidence · `01M6KHKF60GFWGKGN82BPKW8FW` · actor [claude-code-bdf6678e](../actor/claude-code-bdf6678e.md)

## result

PASS

## statement

At 98a46ce3, cmd/rhiza/testdata/script/workflow-new.txtar forks ship-change with --from and checks it (only the name line differs, a fork of the fork is the same bytes, workflow check passes). It writes a skeleton that passes workflow check and run start --dry-run. It now asserts exit 3 on an existing file (with and without --from, plus code conflict under --json) with the file left byte-for-byte unchanged, asserts exit 2 for sub/x and x.toml, and checks --dry-run writes nothing. `go test ./cmd/rhiza -run 'TestScripts/(workflow-new|describe-golden)$'` and `go test ./internal/api -run WithName` pass. By hand, with a scratch binary in a scratch repo (HOME/XDG in the scratchpad, RHIZA_* unset): the fork diffs from internal/workflow/builtin/ship-change.toml only in name and is valid. The skeleton is commented and valid. Repeating new skel exits 3 (also with --from mine) and the sha1 is unchanged. a/b and a.toml exit 2. --dry-run d prints 'would write' and creates no file. With RHIZA_ROOT/CHANGE/RUN/ACTOR unset, fmt-check, vet, self-check and make test pass. One earlier make test run failed under machine-wide lock contention; its output was truncated and a full rerun passed.

## Links

### Out

- verifies [workflow-new](../capability/rhiza/workflows/workflow-authoring/workflow-new.md)

## History

- `01M6HN24SPSS76PKZHZV188JEW` Workflow authoring: rhiza workflow ls / describe / new <name> [--from <builtin>] (forks a built-in into .rhiza/workflows/<name>.toml) / check <file>, plus a guide topic; no workflow ops exist in the registry today (internal/workflow, internal/api run.go workflowResolver). Carve leaves workflow-authoring and workflow-overrides out of umbrella 01M5SWT91P766X (Area workflows)
