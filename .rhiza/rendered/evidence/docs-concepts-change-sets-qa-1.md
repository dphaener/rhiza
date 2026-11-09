# docs-concepts-change-sets-qa-1

evidence · `01M6ZEA7YXC860PB2F4W34PAKH` · actor [claude-code-ca72ee86](../actor/claude-code-ca72ee86.md)

## result

PASS

## statement

QA cycle 1, scratch build in a fresh scratch repo: ran change new (intention and --kind record), change use, --record write, propose, accept --reason --rule, reject, withdraw (writes refused afterwards), change ls/show, land, log, render then check (ok). depends-on between change sets: a loop is refused and land stops naming the draft dependency. carve access-roles --from/--to --with-children moved the node and its child with entries on both sides. reopen took a branch-landed change set back to accepted. All as documented. Small wording mismatch: rhiza log prints 'Carved from X into Y' on both sides rather than 'carved from'/'carved into'.

## Links

### Out

- verifies [docs-concepts-change-sets](../capability/rhiza/documentation/docs-concepts-change-sets.md)

## History

- `01M6W5EM5J6KCHFB3SKS29JJ7M` Concepts docs: docs/ explains the model (node kinds, links, handles, references), change sets and the log, and decisions, sources and why
