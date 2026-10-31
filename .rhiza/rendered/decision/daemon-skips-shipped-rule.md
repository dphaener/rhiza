# daemon-skips-shipped-rule

decision · `01M69DGWF9S6W7NC58Z2AYE839`

## statement

The daemon skips an accepted change set whose latest run is a succeeded ship-change run that ended after the change set was last accepted or reopened; a succeeded plan or review run does not count, and a reopen makes it pickable again; the workflow name is passed into the daemon package, which does not import api

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m8-orchestrator-picks](../source/m8-orchestrator-picks.md)

### In

- decided-by [daemon-skips-shipped](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-skips-shipped.md)
- decided-by [shipped-skip-reason](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-skip-reasons/shipped-skip-reason.md)

## History

- `01M69DFFE95XDD58TYREQR0HXE` Plan milestone 8: Rhiza notices conflicts before they bite, picked by the orchestrator under the owner's standing instruction
