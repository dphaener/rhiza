# question-95eb7493-1

question · `01M6YKW68VV73QQHSM7XMYNJK6`

## session

95eb7493-c225-410a-b431-01a1bd3b689e

## statement

`rhiza change show 01M6WMCCRC0KSCWEV17VY8KVQP` fails because model.ShowFieldValue (internal/model/textshow.go) re-encodes a text field that holds an @@{ID} reference with json.Marshal. That HTML-escapes < > & (and always escapes U+2028/U+2029), so op.Encode's canonical check rejects the shown op. The same helper feeds change show, its indexed path, change diff, show, slice and render. How should the fix work?

## Links

### In

- answers [shown-text-canonical-encoding](../decision/shown-text-canonical-encoding.md)

## History

- `01M6YKP1TFR204NPCNJVQ36A4F` Fix rhiza change show 01M6WMCCRC0KSCWEV17VY8KVQP fails: 'op 01M6X6KYJY5EPC5BJT5R4EEBB6: op: invalid: field "fields": value of "statement": canon: not canonical', while rhiza check passes and rhiza land succeeded on the same change set (PR #383, merged 9c0d4a2a). Reproduces with main build 6f44a0e5 and later. The op's evidence statement contains '<branch>' text; suspected cause: change show's decode path re-encodes the field with encoding/json default HTML escaping (< > &), so the canonical check fails on text that is canonical on disk. Expect: change show (and every read path) accepts any op that check and land accept.
