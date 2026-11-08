# wait-settles-conflict-rule

decision · `01M6A0MCZRX2JM4BS1RDRXK34R`

## statement

ship check --settled reads the PR's mergeable state; on CONFLICTING it runs the ship rebase logic (which returns early when the branch already holds origin/main, so the later rebase steps stay harmless), settling rendered-only conflicts and recording the pushed head, then reports not settled; on a code conflict it raises the rebase-conflict question and exits 75 with a stop: line; wait steps learn to stop on 75 as command steps do; --settled may now write the ship record

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m10-orchestrator-picks](../source/m10-orchestrator-picks.md)

### In

- decided-by [wait-stops-on-75](../capability/rhiza/workflows/workflow-definitions/wait-step/wait-stops-on-75.md)
- decided-by wait-settles-conflicting-pr (retired)

## History

- `01M6A0JENBMBYPYRN5195AGPKH` Plan milestone 10: runs heal themselves, picked by the orchestrator under the owner's standing instruction
