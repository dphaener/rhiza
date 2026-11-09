# shown-text-canonical-encoding

decision · `01M6YKY44M04D82GEPRW690QR9`

## reasoning

Read internal/model/textshow.go, internal/canon/canon.go and op.Encode. ShowFieldValue is the only place the shown text is re-encoded, and it uses json.Marshal, whose default HTML and line-separator escaping differs from canon's emitString; canon.Marshal already exists for exactly this (SetEscapeHTML(false) then Canonicalize) and is what op.Encode uses. Fixing the shared helper covers all six callers named in the question and keeps the canonical check strict, which the change set's intent (every read path accepts what check and land accept) needs; loosening the check or patching callers one by one would leave other paths broken or weaken log integrity.

## statement

Fix it at the one shared helper: model.ShowFieldValue encodes the shown string with canon.Marshal instead of json.Marshal, so a shown text value is always canonical JSON (no HTML escaping of < > &, no forced U+2028/U+2029 escapes) and passes op.Encode's canonical check. Every read path that uses the helper (change show and its indexed path, change diff, show, slice, render) is fixed by that one change, with no per-caller workaround and no relaxing of the canonical check; a regression test covers a referenced text field holding < > & and U+2028, plus a testscript that runs change show on such an op.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-95eb7493-1](../question/question-95eb7493-1.md)
- derived-from [change-show-canon-plan](../source/change-show-canon-plan.md)

### In

- decided-by [shown-text-canonical](../capability/rhiza/model/text-references/shown-text-canonical.md)

## History

- `01M6YKP1TFR204NPCNJVQ36A4F` Fix rhiza change show 01M6WMCCRC0KSCWEV17VY8KVQP fails: 'op 01M6X6KYJY5EPC5BJT5R4EEBB6: op: invalid: field "fields": value of "statement": canon: not canonical', while rhiza check passes and rhiza land succeeded on the same change set (PR #383, merged 9c0d4a2a). Reproduces with main build 6f44a0e5 and later. The op's evidence statement contains '<branch>' text; suspected cause: change show's decode path re-encodes the field with encoding/json default HTML escaping (< > &), so the canonical check fails on text that is canonical on disk. Expect: change show (and every read path) accepts any op that check and land accept.
