# text-ref-check-qa-1

evidence · `01M6KP0VJMY5BBP7XE2TJCC6XG` · actor [claude-code-5c694237](../actor/claude-code-5c694237.md)

## result

FAIL

## statement

QA cycle 1 at 0e9f6a6e, scratch build in a scratch repo: a hand-forged @@{01ZZ…(26 chars)} in a log entry made rhiza check exit 1 with 'fail text-ref ... legacy field statement holds @@{01ZZ...}, which names no node' (and rule text-ref in --json); a reference to a retired node gave a warn while check stayed ok; text-ref appears in describe --json; go test ./internal/check passes. FAIL: rhiza check on this repo at the branch head (make self-check) fails with 'fail text-ref change set 01M6HNNWFR1J2AHQRNKFTPME46: text-ref-check-verified field statement holds @@{01ZZ…(26 chars)}, which names no node'; in a scratch repo, rhiza add with 'a forged @@@{01ZZ…(26 chars)} here' succeeds and then rhiza check fails the same way.

## Links

### Out

- verifies [text-ref-check](../capability/rhiza/model/text-references/text-ref-check.md)

## History

- `01M6HNNWFR1J2AHQRNKFTPME46` Text references: references to other nodes inside statement and other text fields are written as @handle, resolved at write time and stored by ID, and shown by current handle in render, show, why and the read apps; rhiza check flags dangling references (internal/op, internal/render, internal/api/write_node.go). Carve leaf text-references out of umbrella 01M5SWSWDVB237 (Area model)
