# import-guide-topic-qa-1

evidence · `01M6RPTZX7EKDYE4S6FG39T2WX` · actor [claude-code-99a4b355](../actor/claude-code-99a4b355.md)

## result

PASS

## statement

At a909311c, built ./cmd/rhiza to a scratch binary and imported a fixture Go repo (store/, api/ with tests, PLAN.md, DECISIONS.md) by hand, following rhiza guide import. rhiza guide ls lists '- import topic 791 tokens'; rhiza guide import prints steps 1. Init through 7. Verify in order, and each step's commands worked as written (init, apply --record, add decision/link derived-from, source add --raw, change new --source, check, render). go test ./internal/guide/... (TestCaps) and TestScripts/guide pass.

## Links

### Out

- verifies [import-guide-topic](../capability/rhiza/guide/import-guide-topic.md)

## History

- `01M6J36X6EC3E70H0QCW17GREP` Import procedure: rhiza guide import teaches agents to bring an existing repo with built code into Rhiza, the steps Nebari's import follows
