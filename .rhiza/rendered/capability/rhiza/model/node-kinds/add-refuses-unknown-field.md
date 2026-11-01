# add-refuses-unknown-field

capability · `01M69T0JB213041HC3VS3HRQCS` · parent [node-kinds](../node-kinds.md)

## acceptance

A new field-schema.txtar has add evidence --field bogus=x exit 2 naming the known fields, and a hand-written log entry with an unknown field still passes fold-ok and check; scripts that used made-up field names are changed to real ones; make check passes

## statement

rhiza add and rhiza set refuse a field the node's kind does not have, as a usage error naming the known fields, while nodes already in the log with unknown fields still read and fold

## Links

### Out

- decided-by [unknown-fields-refused-on-new-writes](../../../../decision/unknown-fields-refused-on-new-writes.md)

### In

- verifies [add-refuses-unknown-field-verified](../../../../evidence/add-refuses-unknown-field-verified.md)
- verifies [add-refuses-unknown-field-qa-1](../../../../evidence/add-refuses-unknown-field-qa-1.md)
- verifies [pr-146-each-core-node-kind-has-a](../../../../evidence/pr-146-each-core-node-kind-has-a.md)

## History

- `01M5SWTH32C73BQET8E0708F9F` Each core node kind has a field schema, so its fields are known and checked.
