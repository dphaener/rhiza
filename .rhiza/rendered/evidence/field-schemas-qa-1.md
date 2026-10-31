# field-schemas-qa-1

evidence · `01M6A44YGXSHJST0GDCMPKAC0K` · actor [claude-code-ab7a1e6c](../actor/claude-code-ab7a1e6c.md)

## result

PASS

## statement

At cacee34c with a fresh build, every field the internal/op schema lists was exercised through the CLI where the leaf names it (intent statement+umbrella, decision statement+notes_quote, evidence statement/result/observed/head/test, capability statement, actor refusal listing agent, email, for, model, name, role, session) and accepted, while bogus on evidence is refused; the refusal messages print the sorted per-kind lists from internal/op/fields.go. make check (which runs TestEveryFieldInRepoLogIsKnown and rhiza check on this repo's own log) passes with RHIZA_* unset.

## Links

### Out

- verifies [field-schemas](../capability/rhiza/model/node-kinds/field-schemas.md)

## History

- `01M5SWTH32C73BQET8E0708F9F` Each core node kind has a field schema, so its fields are known and checked.
