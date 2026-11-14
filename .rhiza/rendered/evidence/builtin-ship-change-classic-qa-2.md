# builtin-ship-change-classic-qa-2

evidence · `01M7DTKPWTJN83516V0KAA0D9Y` · actor [claude-code-93a5cadf](../actor/claude-code-93a5cadf.md)

## result

PASS

## statement

QA at 1b3323a74, scratch build of the branch: workflow ls lists ship-change-classic (builtin); workflow check ship-change-classic is valid; diff of builtin/ship-change-classic.toml against testdata/ship-change-pre-m36.toml is only the name line, and the golden equals origin/main's ship-change.toml; workflow new --from ship-change-classic writes a file that checks valid; TestBuiltinClassicUnchanged passes.

## Links

### Out

- verifies [builtin-ship-change-classic](../capability/rhiza/workflows/built-in-workflows/builtin-ship-change-classic.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
