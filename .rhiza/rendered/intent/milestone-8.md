# milestone-8

intent · `01M69DGT03S4PSMEZCCAM4G1E5`

## milestone

8

## statement

Rhiza notices conflicts before they bite: the daemon never restarts a change set it just shipped, rhiza check and the inbox warn when in-flight change sets overlap, and plan reviews check the plan against the other change sets in flight

## Links

### Out

- decided-by [m8-picked-by-orchestrator](../decision/m8-picked-by-orchestrator.md)

### In

- serves change set `01M69DFHXFG7197RGNYYPWS4A5` (The daemon does not restart a change set whose latest run is a succeeded ship-change run that ended after the change set was last accepted or reopened)
- serves change set `01M5SWTNWB9GEJ5AA8R7K867KB` (Genesis detects conflicting decisions across change sets in flight and surfaces them to the owner.)
- serves change set `01M5SWTMR4G5EFRVRQHSEKVA15` (Genesis reviews each plan against the model and the other runs in flight before work starts, as an agent check on the plan rather than the code.)

## History

- `01M69DFFE95XDD58TYREQR0HXE` Plan milestone 8: Rhiza notices conflicts before they bite, picked by the orchestrator under the owner's standing instruction
