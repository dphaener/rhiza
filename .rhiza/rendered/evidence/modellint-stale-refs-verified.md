# modellint-stale-refs-verified

evidence · `01M7B11BF6PDX1H8JK5DQQJMV1` · actor [claude-code-69a64eaa](../actor/claude-code-69a64eaa.md)

## result

PASS

## statement

At 4acb20acc: re-ran the Checks block (go test -count=1 -trimpath ./internal/modellint/... ok; rhiza ship generate exit 0, 0 changed; rhiza render exit 0, no tracked diff). go vet and gofmt -l clean. Read the diff: only internal/modellint/{refs.go,refs_test.go} plus the builder's actor log entry and its rendering; no git, no internal/check changes, no TODOs. refs_test.go covers each acceptance point: present path/test/script report nothing (TestStaleRefsPresent); missing path, Go test and testscript each report one hit naming node, field and ref (table TestStaleRefsMissing); all-paths-gone capability reported unretired (TestStaleRefsUnretired); an unwritten node with a missing path reports nothing (TestStaleRefsOtherNodesNotRead); and/or and a URL are not references (TestStaleRefsPresent, TestStaleRefsNotReferences).

## Links

### Out

- verifies [modellint-stale-refs](../capability/rhiza/workflows/built-in-workflows/built-in-review/modellint-stale-refs.md)

## History

- `01M76XTF51R02XC5WXJSSBP8WP` Review's code step runs deterministic model lints before the model check: acceptance naming tests or files that no longer exist, Decisions never applied, capabilities whose code is gone but not retired. Evidence: the model check fails cycle 1 in 37-42% of runs, mostly mechanical drift (perf report rec 4).
