# shown-text-canonical-verified

evidence · `01M6Z0PDG1C8VYCDN048XSAAV1` · actor [claude-code-a4d2d949](../actor/claude-code-a4d2d949.md)

## result

PASS

## statement

Verifier claude-code-a4d2d949 at 5596c942: read the diff (ShowFieldValue now encodes with canon.Marshal minus the trailing newline). go test -trimpath ./internal/model -run ShowField passes; TestShowFieldValueCanonical feeds a statement with an ID reference with <branch>, a & b, x > y and U+2028 (bytes checked), asserts equality with canon.Marshal of the shown text minus newline and canon.Check. go test -trimpath ./cmd/rhiza -run TestScripts/change-show-canonical-text passes (with GOTMPDIR and the run's RHIZA_* env unset; under the run env every script, incl. existing change-show-serves, fails at rhiza init by resolving /home/rhiza/code/rhiza): check, change show, change show --json, change diff, show mine --json print literal < > &, and change show of the landed change set succeeds. With the two tests on the parent commit both fail: canon not canonical in the unit test, and rhiza check's next step change show rejects the op as canon: not canonical. Whole internal/model package and change-show*/change-diff*/show* scripts pass.

## Links

### Out

- verifies [shown-text-canonical](../capability/rhiza/model/text-references/shown-text-canonical.md)

## History

- `01M6YKP1TFR204NPCNJVQ36A4F` Fix rhiza change show 01M6WMCCRC0KSCWEV17VY8KVQP fails: 'op 01M6X6KYJY5EPC5BJT5R4EEBB6: op: invalid: field "fields": value of "statement": canon: not canonical', while rhiza check passes and rhiza land succeeded on the same change set (PR #383, merged 9c0d4a2a). Reproduces with main build 6f44a0e5 and later. The op's evidence statement contains '<branch>' text; suspected cause: change show's decode path re-encodes the field with encoding/json default HTML escaping (< > &), so the canonical check fails on text that is canonical on disk. Expect: change show (and every read path) accepts any op that check and land accept.
