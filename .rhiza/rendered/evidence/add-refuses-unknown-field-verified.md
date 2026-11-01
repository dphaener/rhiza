# add-refuses-unknown-field-verified

evidence · `01M6A3N7TADFYRBFMVR71PAK04` · actor [claude-code-234feb1c](../actor/claude-code-234feb1c.md)

## command

go test ./cmd/rhiza -run TestScripts/field-schema -count=1 -v && make check

## result

PASS

## statement

At 776785fc: new cmd/rhiza/testdata/script/field-schema.txtar passes (go test -run TestScripts/field-schema): add evidence --field bogus=x is refused naming the known fields (command, head, observed, pr, result, statement, status, test, tests), --json gives code usage, anchor says it has no fields, set refuses on a new and a landed node, nothing is written by refusals, and a log entry hand-edited to an unknown field (sed command->bogus) passes fold-ok, render and check. By hand in a scratch repo with a fresh build, add evidence e1 --field bogus=x and set e1 --field bogus=x both exit 2. Scripts that used made-up fields (tags, weight, owner, a/b/x/s, actor on evidence, milestone on capability) now use real kind fields; the one left (record.txtar set draft --field a=b) still fails on not_found first as intended. Diff: knownFieldsOnly in internal/api/write_node.go guards buildAdd and buildSet only, so folding is untouched; notes_quote added to Decision fields. make check exits 0 with the run RHIZA_* env unset (with it set, self-check only flags the verifier's own actor render from RHIZA_ROOT).

## Links

### Out

- verifies [add-refuses-unknown-field](../capability/rhiza/model/node-kinds/add-refuses-unknown-field.md)

## History

- `01M5SWTH32C73BQET8E0708F9F` Each core node kind has a field schema, so its fields are known and checked.
