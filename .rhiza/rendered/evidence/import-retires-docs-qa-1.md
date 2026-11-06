# import-retires-docs-qa-1

evidence · `01M6RPV03WWCFE5SM1BJBAJCVX` · actor [claude-code-99a4b355](../actor/claude-code-99a4b355.md)

## result

PASS

## statement

At a909311c, built ./cmd/rhiza to a scratch binary and imported a fixture Go repo (store/, api/ with tests, PLAN.md, DECISIONS.md) by hand, following rhiza guide import. rhiza source add decisions-md/plan-md <file> --raw <doc> succeeded with raw_sha256 shown by rhiza why; git rm of DECISIONS.md and PLAN.md went in the same commit as the Source's log records; the change set landed and rhiza check printed ok.

## Links

### Out

- verifies [import-retires-docs](../capability/rhiza/guide/import-retires-docs.md)

## History

- `01M6J36X6EC3E70H0QCW17GREP` Import procedure: rhiza guide import teaches agents to bring an existing repo with built code into Rhiza, the steps Nebari's import follows
