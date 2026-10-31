# check-warns-in-flight-link-kinds-verified-2

evidence · `01M6896K9XQH0JSWXZGW227GQN` · actor [claude-code-35812831](../actor/claude-code-35812831.md)

## result

PASS

## statement

Verifier attempt 2 at 49b80f8. go test ./cmd/rhiza -run TestScripts/check passes; check.txtar now swaps the from/to IDs as well as the type, so the in-flight entry holds a verifies link from intent fast to capability store in an open (draft) change set, asserts rhiza check exits 0 printing 'warn in-flight change set <id>, op <id>, link <id>: link type "verifies" does not go from intent to capability', and asserts no in-flight line after rhiza withdraw. Replayed the same steps by hand with a scratch build: the rewritten entry's from is the intent's ID, check exits 0 with exactly that warning, and the withdrawn change set gives none. make fmt-check vet test all pass. make self-check fails only on render staleness for the builder's own actor file claude-code-94bcc010.md and index.md, whose actor log entry the run records outside this commit; nothing to do with the leaf's code.

## Links

### Out

- verifies [check-warns-in-flight-link-kinds](../capability/rhiza/integrity/check/check-warns-in-flight-link-kinds.md)

## History

- `01M5SWTV1D5DNN6E2B7HB1XA9Q` genesis link checks endpoint kinds when the link is written, not only in genesis check.
