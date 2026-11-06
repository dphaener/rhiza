# import-unbuilt-as-drafts-verified

evidence · `01M6RNM0BV5KFGHV73TYYV8ZAD` · actor [claude-code-483c1206](../actor/claude-code-483c1206.md)

## result

PASS

## statement

Verifier at 554595ee: read the diff (internal/guide/text/import.md step 6 now says to add the planning document as a Source, open a draft change set per unbuilt piece with rhiza change new "<intent>" --source <source>, leave it in draft since accepting is the owner's, and write each open question as a question node via rhiza add question). Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(guide|import)' and go test -trimpath ./internal/guide/...: all PASS, including the new import-unbuilt-as-drafts.txtar, which imports a fixture PLAN.md with one TODO and one TBD and ends with one draft intention change set whose Sources names the plan Source, one question node export-format, no accepted change set, and fold-ok. Reproduced by hand with a scratch build in a temp repo: change show reported State: draft, Sources: <source id>; change ls listed only landed/draft change sets.

## Links

### Out

- verifies [import-unbuilt-as-drafts](../capability/rhiza/guide/import-unbuilt-as-drafts.md)

## History

- `01M6J36X6EC3E70H0QCW17GREP` Import procedure: rhiza guide import teaches agents to bring an existing repo with built code into Rhiza, the steps Nebari's import follows
