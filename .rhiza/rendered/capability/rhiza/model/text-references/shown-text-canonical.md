# shown-text-canonical

capability · `01M6YM01FKQKHVJ1VCA63B6T2J` · parent [text-references](../text-references.md)

## acceptance

A unit test in internal/model (go test -trimpath ./internal/model -run ShowField) gives ShowFieldValue a statement with an @@{ID} reference plus <branch>, a & b, x > y and U+2028, and asserts the result equals canon.Marshal of the shown text minus its newline and passes canon.Check; a testscript (go test -trimpath ./cmd/rhiza -run TestScripts/change-show-canonical-text) writes such a statement referring to another node, then rhiza check passes, rhiza change show and change show --json succeed and print the text with literal <, > and &, rhiza change diff and rhiza show --json of the node succeed with the same literal text, and after landing, change show of the landed change set still succeeds

## statement

A text field shown by handle stays canonical JSON (literal < > & and U+2028/U+2029, as canon writes them), so change show, its indexed path, change diff, show, slice and render read every op that rhiza check and rhiza land accept

## Links

### Out

- decided-by [shown-text-canonical-encoding](../../../../decision/shown-text-canonical-encoding.md)
- derived-from [change-show-canon-plan](../../../../source/change-show-canon-plan.md)

### In

- verifies [shown-text-canonical-verified](../../../../evidence/shown-text-canonical-verified.md)
- verifies [shown-text-canonical-qa-1](../../../../evidence/shown-text-canonical-qa-1.md)
- verifies [pr-404-fix-rhiza-change-show](../../../../evidence/pr-404-fix-rhiza-change-show.md)

## History

- `01M6YKP1TFR204NPCNJVQ36A4F` Fix rhiza change show 01M6WMCCRC0KSCWEV17VY8KVQP fails: 'op 01M6X6KYJY5EPC5BJT5R4EEBB6: op: invalid: field "fields": value of "statement": canon: not canonical', while rhiza check passes and rhiza land succeeded on the same change set (PR #383, merged 9c0d4a2a). Reproduces with main build 6f44a0e5 and later. The op's evidence statement contains '<branch>' text; suspected cause: change show's decode path re-encodes the field with encoding/json default HTML escaping (< > &), so the canonical check fails on text that is canonical on disk. Expect: change show (and every read path) accepts any op that check and land accept.
