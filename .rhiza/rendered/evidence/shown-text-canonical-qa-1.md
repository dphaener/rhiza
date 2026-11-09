# shown-text-canonical-qa-1

evidence · `01M6Z15YM0ASWDCAYMVE3X711E` · actor [claude-code-97ba9ca3](../actor/claude-code-97ba9ca3.md)

## result

PASS

## statement

QA at 7c4c8dad with the built binary. On the real repo, the merge-base build fails rhiza change show 01M6WMCCRC0KSCWEV17VY8KVQP (op 01M6X6KYJY5EPC5BJT5R4EEBB6: canon: not canonical); the branch build prints it, text and --json, with literal <branch>. In a scratch repo, a capability statement with @@other, <branch>, a & b, x > y and U+2028/U+2029: check passes; change show, change show --json (no u003c or u2028 escapes), change diff, show --json, slice and render succeed with literal text; after propose/accept/land, change show of the landed change set succeeds (the base build fails it as not canonical). go test -trimpath ./internal/model -run ShowField and ./cmd/rhiza -run TestScripts/change-show-canonical-text: ok.

## Links

### Out

- verifies [shown-text-canonical](../capability/rhiza/model/text-references/shown-text-canonical.md)

## History

- `01M6YKP1TFR204NPCNJVQ36A4F` Fix rhiza change show 01M6WMCCRC0KSCWEV17VY8KVQP fails: 'op 01M6X6KYJY5EPC5BJT5R4EEBB6: op: invalid: field "fields": value of "statement": canon: not canonical', while rhiza check passes and rhiza land succeeded on the same change set (PR #383, merged 9c0d4a2a). Reproduces with main build 6f44a0e5 and later. The op's evidence statement contains '<branch>' text; suspected cause: change show's decode path re-encodes the field with encoding/json default HTML escaping (< > &), so the canonical check fails on text that is canonical on disk. Expect: change show (and every read path) accepts any op that check and land accept.
