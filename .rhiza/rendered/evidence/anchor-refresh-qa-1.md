# anchor-refresh-qa-1

evidence · `01M6JAFA2A7N2F82S09V5BZM73` · actor [claude-code-666954d5](../actor/claude-code-666954d5.md)

## result

PASS

## statement

Built the branch head (9a4dfa5d) into a scratch binary and drove it by hand in a scratch git repo with RHIZA_ROOT/CHANGE/RUN unset. In a change set, anchor refresh pin --dry-run --json and api call anchor.refresh --dry-run printed identical bytes, and the dry runs left the .rhiza file count unchanged. A real refresh wrote the new fingerprint (previous to new), and a second refresh said it was unchanged and wrote nothing. A gone symbol and a gone path were each refused with exit 2 and a next line. After propose, accept, land and render, check had no finding for pin. The anchor-refresh, describe-golden and other anchor scripts pass.

## Links

### Out

- verifies [anchor-refresh](../capability/rhiza/integrity/fingerprints/anchor-refresh.md)

## History

- `01M6HJXAJ6GYA2RTAGVK6JR1N9` Code anchors with fingerprints: an anchor node carries path, optional symbol and a content fingerprint; implements links it to a capability; rhiza check resolves every anchor and warns when the path or symbol is gone or the fingerprint changed (internal/op/fields.go says anchors carry no fields yet). Carve leaves anchors and fingerprints out of umbrella 01M5SWT0K9TV1R (Area integrity) rather than duplicating them
