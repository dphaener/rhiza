# release-gates-verified

evidence · `01M785YSX1PVZCD0JR7KZ5FE07` · actor [claude-code-988bdbad](../actor/claude-code-988bdbad.md)

## result

PASS

## statement

Verifier at 4bf1b104d (attempt 2): ran go test -trimpath -count=1 ./tools/release/... -run TestGates -v; TestGates (14 subtests: each gate passing and failing, commit missing from history fails soak, severe skips soak and batching but never CI) and TestGatesCheckReports pass. TestGates now reads its history from the fixture testdata/runs.jsonl through ReadRunHistory (the JSON-lines format tools/release defines), fixing attempt 1's f2. tools/release check now wires the gates into runCheck (main.go parses -runs/RHIZA_RELEASE_RUNS and -commit; CI via an injectable gh lookup, bug requests and blocking questions from the log, last patch tag from git) and prints each gate with its reason, fixing f1; TestGatesCheckReports drives runCheck over a fixture repo with a fake clock, fake CI and fixture history and checks ready/blocked output, CI red, soak after a bug request, missing commit, and a blocking question cleared by an answer. Full package tests, go vet and gofmt clean.

## Links

### Out

- verifies [release-gates](../capability/rhiza/distribution/release-gates.md)

## History

- `01M76B0FZA3N7ZHAYTBN76MVSV` rhiza release check: read-only; computes the next version and its kind from landed milestones, fixes and breaking marks or detections, evaluates the gates (CI green, 24h soak of the latest build, milestone closed, severe-fix skip) and drafts the notes from the model
