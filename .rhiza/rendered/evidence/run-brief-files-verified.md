# run-brief-files-verified

evidence · `01M7AQWK34TXZACQVJB05T0GMD` · actor [claude-code-0b340a5c](../actor/claude-code-0b340a5c.md)

## result

PASS

## statement

At 92d4ec944: re-ran the Checks: block and acceptance. go test -trimpath ./internal/run -run 'TestBriefFiles|TestSessionDir|TestHandbackFile' exit 0 (TestBriefFiles, TestSessionDir, TestSessionDirWrite, TestHandbackFile pass); rhiza ship generate exit 0 (0 files changed); rhiza render exit 0 (only this verifier's own actor written); gofmt -l and go vet ./internal/run clean. Read the diff: only internal/run/brief.go and brief_test.go plus the builder's actor log/rendered files; internal/api and the workflow parser untouched; no TODOs. brief.md sections in order (change set title/ID/intent, leaves with acceptance, decisions, constraints, check commands with 'none' when neither set, contract ending 'Full contract: rhiza guide contract') and byte-identical re-render are tested; TestSessionDir checks all five session-name forms incl. leaf fixer/verifier and change set fixer (answering review question r3-f1) and asserts all distinct; handback.md tested for gate commands with output paths, CI runs with log files, conflicted files; failure cases tested: uncreatable directory error names its path, cause with no step refused. WriteHandback has no callers outside its package, as required.

## Links

### Out

- verifies [run-brief-files](../capability/rhiza/workflows/workflow-runs/run-brief-files.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
