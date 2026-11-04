# check-fingerprint-changed-qa-1

evidence · `01M6JAEZQMDJ1WK8WBTVF90VG9` · actor [claude-code-666954d5](../actor/claude-code-666954d5.md)

## result

PASS

## statement

Built the branch head (9a4dfa5d) into a scratch binary and drove it by hand in a scratch git repo with RHIZA_ROOT/CHANGE/RUN unset. After the anchored function's body changed, check exited 0 with one warn naming anchor pin and capability thing, and with next 'rhiza anchor refresh pin'. Editing only the doc comment above the declaration added no finding. Appending to a whole-file-anchored notes.md gave a 'file notes.md changed' warn. go test ./internal/check passes.

## Links

### Out

- verifies [check-fingerprint-changed](../capability/rhiza/integrity/fingerprints/check-fingerprint-changed.md)

## History

- `01M6HJXAJ6GYA2RTAGVK6JR1N9` Code anchors with fingerprints: an anchor node carries path, optional symbol and a content fingerprint; implements links it to a capability; rhiza check resolves every anchor and warns when the path or symbol is gone or the fingerprint changed (internal/op/fields.go says anchors carry no fields yet). Carve leaves anchors and fingerprints out of umbrella 01M5SWT0K9TV1R (Area integrity) rather than duplicating them
