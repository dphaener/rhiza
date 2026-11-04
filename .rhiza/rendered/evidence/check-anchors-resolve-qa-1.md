# check-anchors-resolve-qa-1

evidence · `01M6JAEHDG1C5YQS3N11DQW47A` · actor [claude-code-666954d5](../actor/claude-code-666954d5.md)

## result

PASS

## statement

Built the branch head (9a4dfa5d) into a scratch binary and drove it by hand in a scratch git repo with RHIZA_ROOT/CHANGE/RUN unset. With a valid anchor, check printed ok. After renaming the anchored func Keep, check exited 1 with 'fail anchor pkg/p.go: anchor pin (implements thing): symbol Keep of pkg/p.go is gone'. After moving the file away, it printed a 'path pkg/p.go is gone' fail for each of the three anchors on that file, naming the capability, or 'implements nothing' when unlinked. Once the file was restored, check was clean again.

## Links

### Out

- verifies [check-anchors-resolve](../capability/rhiza/integrity/anchors/check-anchors-resolve.md)

## History

- `01M6HJXAJ6GYA2RTAGVK6JR1N9` Code anchors with fingerprints: an anchor node carries path, optional symbol and a content fingerprint; implements links it to a capability; rhiza check resolves every anchor and warns when the path or symbol is gone or the fingerprint changed (internal/op/fields.go says anchors carry no fields yet). Carve leaves anchors and fingerprints out of umbrella 01M5SWT0K9TV1R (Area integrity) rather than duplicating them
