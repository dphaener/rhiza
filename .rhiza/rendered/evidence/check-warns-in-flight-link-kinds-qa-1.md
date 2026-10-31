# check-warns-in-flight-link-kinds-qa-1

evidence · `01M687RCCQRJDW6B2NR5280TBA` · actor [claude-code-d48993a5](../actor/claude-code-d48993a5.md)

## result

PASS

## statement

QA at 9c8d4bf with a scratch build in a fresh repo: in an open change set I wrote link store serves fast, then edited the log entry into a verifies link from intent fast to capability store; rhiza check exits 0, prints ok and warn in-flight change set <id>, op <id>, link <id>: link type "verifies" does not go from intent to capability, both while draft and after propose; after rhiza withdraw, check exits 0 with no in-flight line. In a second repo the same forged link warns in-flight once while accepted, and after rhiza land it becomes one fail link line with no in-flight warning.

## Links

### Out

- verifies [check-warns-in-flight-link-kinds](../capability/rhiza/integrity/check/check-warns-in-flight-link-kinds.md)

## History

- `01M5SWTV1D5DNN6E2B7HB1XA9Q` genesis link checks endpoint kinds when the link is written, not only in genesis check.
