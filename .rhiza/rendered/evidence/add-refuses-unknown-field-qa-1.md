# add-refuses-unknown-field-qa-1

evidence · `01M6A44X44XGEMQKW0V0Z4J4GZ` · actor [claude-code-ab7a1e6c](../actor/claude-code-ab7a1e6c.md)

## result

PASS

## statement

At cacee34c, built go build -o $TMPDIR/rhiza-qa ./cmd/rhiza and used it in a scratch git repo with the run env unset. add evidence e1 --field bogus=x exits 2 with 'an evidence has no field "bogus" (known fields: command, head, observed, pr, result, statement, status, test, tests)', next 'rhiza describe add'; --json gives code usage; capability with zeta and alpha names alpha first; anchor says it has no fields; actor --field nickname refused; rhiza apply with weight=3 refused as writes[0] usage; api call add refused the same; set e2 --field bogus=x exits 2 naming known fields; nothing refused shows in change show. Known fields (intent umbrella, decision notes_quote, evidence observed/head) are accepted. After hand-editing the log so e2 carries an unknown field legacyx, show e2 lists it, set e2 --field test=t1 still succeeds, render and check pass. set --unset of an unknown name is allowed. make check passes on this checkout with RHIZA_* unset.

## Links

### Out

- verifies [add-refuses-unknown-field](../capability/rhiza/model/node-kinds/add-refuses-unknown-field.md)

## History

- `01M5SWTH32C73BQET8E0708F9F` Each core node kind has a field schema, so its fields are known and checked.
