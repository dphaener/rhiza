# import-unbuilt-as-drafts-qa-1

evidence · `01M6RPV06FRVFEED44H77EC5SP` · actor [claude-code-99a4b355](../actor/claude-code-99a4b355.md)

## result

PASS

## statement

At a909311c, built ./cmd/rhiza to a scratch binary and imported a fixture Go repo (store/, api/ with tests, PLAN.md, DECISIONS.md) by hand, following rhiza guide import. rhiza change new 'Add an export command' --source plan-md left a draft intention change set whose show lists the PLAN.md Source; the TBD became question export-format; nothing was accepted by the import step, and change ls still shows it draft after check.

## Links

### Out

- verifies [import-unbuilt-as-drafts](../capability/rhiza/guide/import-unbuilt-as-drafts.md)

## History

- `01M6J36X6EC3E70H0QCW17GREP` Import procedure: rhiza guide import teaches agents to bring an existing repo with built code into Rhiza, the steps Nebari's import follows
