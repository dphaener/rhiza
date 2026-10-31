# field-schemas

capability · `01M69T0H3RSV8MEH7JNM194ZG4` · parent [node-kinds](../node-kinds.md)

## acceptance

A Go test in internal/op asserts every field used in this repo's log is known for its kind (intent statement, milestone, umbrella; capability statement, acceptance, built_in; evidence statement, result, tests, test, status, head, pr, command, observed; and the decision, source, question, constraint and actor fields) and that bogus on an evidence node is unknown

## statement

Each core node kind has a fixed list of known fields in internal/op, covering every field in use today

## Links

### Out

- decided-by [unknown-fields-refused-on-new-writes](../../../../decision/unknown-fields-refused-on-new-writes.md)

### In

- verifies [field-schemas-verified](../../../../evidence/field-schemas-verified.md)
- verifies [field-schemas-qa-1](../../../../evidence/field-schemas-qa-1.md)
- verifies [pr-146-each-core-node-kind-has-a](../../../../evidence/pr-146-each-core-node-kind-has-a.md)

## History

- `01M5SWTH32C73BQET8E0708F9F` Each core node kind has a field schema, so its fields are known and checked.
