# check-warns-in-flight-link-kinds-verified

evidence · `01M687MSHTBCBD6QZ3TBPK0NJM` · actor [claude-code-367c2fd4](../actor/claude-code-367c2fd4.md)

## result

PASS

## statement

At 27a2f6d: go test ./... passes, including TestScripts/check (check.txtar's new in-flight case: a sed-edited log entry turns a serves link into verifies; rhiza check exits 0, prints ok and warn in-flight change set <id>, op <id>, link <id>: link type "verifies" does not go from capability to intent, no fail; after rhiza withdraw no in-flight line). The script's link runs capability to intent rather than the acceptance's intent to capability, so I reproduced the acceptance by hand with a scratch build in a fresh repo: in an open change set I edited the link.create log entry to a verifies link from intent fast to capability store; rhiza check exited 0 and printed warn in-flight change set 01M687M4116DXQX5K6E4HFVD9E, op 01M687M44HFHT2YT6CS6QTCYPM, link 01M687M44HFHT2YT6CS3M2J693: link type "verifies" does not go from intent to capability; after rhiza withdraw, check exited 0 with no in-flight warning. Diff read: inFlightLinks in internal/check/check.go warns only for draft, proposed and accepted sets, on links that set created in its what-if, each once, sharing badLinks/linkFinding with the fail-severity link rule. make check's fmt-check, vet and test pass; its self-check fails only on render staleness from this verifier session's own actor record, not on the commit.

## Links

### Out

- verifies [check-warns-in-flight-link-kinds](../capability/rhiza/integrity/check/check-warns-in-flight-link-kinds.md)

## History

- `01M5SWTV1D5DNN6E2B7HB1XA9Q` genesis link checks endpoint kinds when the link is written, not only in genesis check.
