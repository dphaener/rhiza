# import-records-built-verified

evidence · `01M6RMMC7HVN7QVQ47FD8NDKCM` · actor [claude-code-e27e0411](../actor/claude-code-e27e0411.md)

## result

PASS

## statement

Verifier at 53bbe72e: read the diff. internal/guide/text/import.md step 3 now states one record per area (never one for the whole repo), Statement: one user-visible behaviour per capability, Acceptance: the repo's existing tests named by test file (a question instead of an invented acceptance), and a new 'Follow-up: anchors' section with rhiza add anchor and link implements, then rhiza check. Ran go test -trimpath ./cmd/rhiza -run TestScripts/import-records-built -count=1 -v: PASS. The script greps those rules from rhiza guide import, records two areas on a fixture repo with Go tests (change ls shows 2 'landed record: add capability'), each capability's acceptance naming its _test.go file, then a follow-up change set adding anchors store-go/api-go linked implements, lands it, renders, and rhiza check prints ok; fold-ok passes. Also go test ./internal/guide/... and ./internal/api -run Guide pass.

## Links

### Out

- verifies [import-records-built](../capability/rhiza/guide/import-records-built.md)

## History

- `01M6J36X6EC3E70H0QCW17GREP` Import procedure: rhiza guide import teaches agents to bring an existing repo with built code into Rhiza, the steps Nebari's import follows
