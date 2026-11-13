# modellint-branch-removed-code-verified

evidence · `01M7B1D1FF2T328ZWJ5Q41TVX8` · actor [claude-code-59fe0de1](../actor/claude-code-59fe0de1.md)

## result

PASS

## statement

At 3fb7c8652 re-ran the Checks: block: go test -trimpath -count=1 ./internal/modellint/... exit 0, rhiza ship generate exit 0 (0 changed), rhiza render exit 0 (only this verifier's own actor written). Read the diff: internal/modellint/removed.go adds RemovedCode(model, change set, fs.FS, deleted []string), which skips retired capabilities and nodes the change set writes, counts path references (first segment in head, or deleted), and reports a capability only when every path reference is absent from head and at least one is in the deleted list (directory refs match deleted paths under them); it runs no git. removed_test.go table tests cover: only path deleted reported naming capability and path; directory deleted reported; deleted path plus a present path not reported; missing but not deleted not reported; retired not reported; no path references not reported; written by the change set not reported. Only internal/modellint plus the actor record and its rendered files changed; no TODOs; gofmt and go vet clean.

## Links

### Out

- verifies [modellint-branch-removed-code](../capability/modellint-branch-removed-code.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
