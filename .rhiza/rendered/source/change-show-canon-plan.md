# change-show-canon-plan

source · `01M6YKZ5R257WA9KNBFH4X3RFD`

## quotes

```json
[]
```

## raw_sha256

91fc6ce94c635d7aa2b391605c71f2b9293c8593e3eab59e9130b934d0fc00ca

## session

95eb7493-c225-410a-b431-01a1bd3b689e

## summary

Autonomous plan run, no owner. Reproduced: rhiza change show 01M6WMCCRC0KSCWEV17VY8KVQP fails on evidence op 01M6X6KYJY5EPC5BJT5R4EEBB6, whose statement holds both an @{ID} text reference and '<branch>'. Cause: model.ShowFieldValue (internal/model/textshow.go) re-encodes a referenced text with json.Marshal, which HTML-escapes < > & and always escapes U+2028/U+2029; op.Encode then rejects the shown op as not canonical. The helper also feeds change show's indexed path, change diff, show, slice and render. One question went to the arbiter (fix shape); it ruled to encode with canon.Marshal in the shared helper, keep op.Encode's check strict, and cover it with a unit test and a change show testscript. The change set serves milestone-25 (backlog and bugs), as the sibling ship-change fix does.

## title

Plan: fix rhiza change show on canonical text (autonomous)

## Links

### In

- derived-from [shown-text-canonical-encoding](../decision/shown-text-canonical-encoding.md)
- derived-from [shown-text-canonical](../capability/rhiza/model/text-references/shown-text-canonical.md)

## History

- `01M6YKP1TFR204NPCNJVQ36A4F` Fix rhiza change show 01M6WMCCRC0KSCWEV17VY8KVQP fails: 'op 01M6X6KYJY5EPC5BJT5R4EEBB6: op: invalid: field "fields": value of "statement": canon: not canonical', while rhiza check passes and rhiza land succeeded on the same change set (PR #383, merged 9c0d4a2a). Reproduces with main build 6f44a0e5 and later. The op's evidence statement contains '<branch>' text; suspected cause: change show's decode path re-encodes the field with encoding/json default HTML escaping (< > &), so the canonical check fails on text that is canonical on disk. Expect: change show (and every read path) accepts any op that check and land accept.
