# check-warns-in-flight-link-kinds

capability · `01M6857QTJSC1V4R4FRDD8ZX3X` · parent [check](../check.md)

## acceptance

check.txtar writes a log entry file holding a verifies link from an intent to a capability in an open change set; rhiza check exits 0 and prints warn in-flight change set <id>, op <id>, link <id>: link type "verifies" does not go from intent to capability; a withdrawn change set holding the same link gives no warning

## statement

rhiza check warns, rule in-flight, on a link with bad endpoint kinds that is only in a draft, proposed or accepted change set, naming the change set, op and link, so a bad link written by hand or by an older binary shows before it lands

## Links

### Out

- decided-by [link-kinds-refused-at-write](../../../../decision/link-kinds-refused-at-write.md)

### In

- verifies [check-warns-in-flight-link-kinds-verified](../../../../evidence/check-warns-in-flight-link-kinds-verified.md)
- verifies [check-warns-in-flight-link-kinds-qa-1](../../../../evidence/check-warns-in-flight-link-kinds-qa-1.md)
- verifies [check-warns-in-flight-link-kinds-verified-2](../../../../evidence/check-warns-in-flight-link-kinds-verified-2.md)
- verifies [pr-124-genesis-link-checks-endpoint](../../../../evidence/pr-124-genesis-link-checks-endpoint.md)

## History

- `01M5SWTV1D5DNN6E2B7HB1XA9Q` genesis link checks endpoint kinds when the link is written, not only in genesis check.
