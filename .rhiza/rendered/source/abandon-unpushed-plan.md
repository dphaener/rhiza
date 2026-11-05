# abandon-unpushed-plan

source · `01M6JX2TS2RMRND0XD3JTHR6Q9`

## quotes

```json
[]
```

## raw_sha256

dbe1852ab716065f1de183a237fdf75cb7782d55f449571ea15afeb46c694d9b

## session

b947d57e-7bee-4636-b577-01a1a5cce3e1

## summary

Autonomous plan run, no owner. Found that PR #164 already switched rhiza run abandon to git branch -d, and run-abandon-unmerged.txtar covers a never-pushed branch with no upstream. The model had no leaf for it, and nothing tested a branch with an upstream plus local commits never pushed to it. Linked the change set serves milestone-25 and relates-to the leafless draft 01M6B9FCNTTH65, then withdrew that draft. Raised the scope question; the arbiter ruled (Decision abandon-unpushed-branch-scope): one childless leaf under run-abandon, a test for the upstream case beside the no-upstream one, and a code change only if that test shows the branch deleted.

## title

Plan: abandon keeps an unpushed branch (autonomous)

## Links

### In

- derived-from [abandon-keeps-unpushed-branch](../capability/rhiza/workflows/workflow-runs/run-abandon/abandon-keeps-unpushed-branch.md)

## History

- `01M6JWSCZYYQVFR6ZTVWS6JWJ0` Fix rhiza run abandon keeps a local branch holding commits that were never pushed, instead of force-deleting them (replaces draft change set 01M6B9FCNTTH65, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
