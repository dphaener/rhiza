# import-retires-docs-verified

evidence · `01M6RPCPS1Q423ZYN3232QRVAP` · actor [claude-code-3ca7c797](../actor/claude-code-3ca7c797.md)

## result

PASS

## statement

At de15d37c, go test -trimpath -count=1 ./cmd/rhiza -run TestScripts/import- passes (import-retires-docs included). Read the diff: rhiza guide import section 5 names retired documents (built or overtaken PLAN, roadmap, TODO, design note: added as a Source with rhiza source add <handle> <file> --raw <document> plus git rm, in the same change set; scrubbed summary and the owner's typed words verbatim) and kept ones (reference, design and operations documents that capabilities link to stay in git); the script asserts those lines. import-retires-docs.txtar commits a fixture PLAN.md, runs rhiza source add plan-md source.json --raw PLAN.md and git rm PLAN.md, commits once, and checks git show --name-status HEAD has both D PLAN.md and A .rhiza/log records, .rhiza/log is clean, the Source carries a sha256, and after propose/accept/land/render rhiza check prints ok and fold-ok passes. This addresses attempt 1's finding f1.

## Links

### Out

- verifies [import-retires-docs](../capability/rhiza/guide/import-retires-docs.md)

## History

- `01M6J36X6EC3E70H0QCW17GREP` Import procedure: rhiza guide import teaches agents to bring an existing repo with built code into Rhiza, the steps Nebari's import follows
