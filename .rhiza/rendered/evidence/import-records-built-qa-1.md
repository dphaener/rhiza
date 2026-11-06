# import-records-built-qa-1

evidence · `01M6RPTZZJZR05T8PGNQF9AM12` · actor [claude-code-99a4b355](../actor/claude-code-99a4b355.md)

## result

PASS

## statement

At a909311c, built ./cmd/rhiza to a scratch binary and imported a fixture Go repo (store/, api/ with tests, PLAN.md, DECISIONS.md) by hand, following rhiza guide import. rhiza apply <file> --record --json per area gave two landed record change sets (store, api), each capability with a statement and acceptance naming its test file; a follow-up change set added anchors store-go/api-go linked implements, proposed/accepted/landed, and rhiza check printed ok. Moving store/store.go made rhiza check fail on anchor store-go, as the topic says.

## Links

### Out

- verifies [import-records-built](../capability/rhiza/guide/import-records-built.md)

## History

- `01M6J36X6EC3E70H0QCW17GREP` Import procedure: rhiza guide import teaches agents to bring an existing repo with built code into Rhiza, the steps Nebari's import follows
