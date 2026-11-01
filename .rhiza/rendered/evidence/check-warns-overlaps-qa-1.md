# check-warns-overlaps-qa-1

evidence · `01M6A9VGC44GQQ9JR3MRCQFNWM` · actor [claude-code-687f5612](../actor/claude-code-687f5612.md)

## result

PASS

## statement

QA at 158cf26b with a scratch build (go build ./cmd/rhiza) in a fresh git repo, RHIZA_* env cleared. Two proposed change sets setting capability thing's statement: rhiza check prints ok, exit 0, and one 'warn in-flight change set A: change sets A and B both write field:statement of <thing id>' with next 'rhiza slice --change A --view what-if'; that slice command runs (exit 0) and shows the proposed edit and the other set in flight; check --json has ok true and the warn under rule in-flight. A draft on the same field adds no overlap warn; a proposed edit of another node adds none. An accepted vs a proposed set on one field warns once. A proposed retire of intent goal vs a proposed link thing serves goal warns 'both write linked, retired of <goal id>'. After accepting and landing A the overlap warn is gone and only the existing stale-base conflict warns remain.

## Links

### Out

- verifies [check-warns-overlaps](../capability/rhiza/integrity/check/check-warns-overlaps.md)

## History

- `01M5SWTNWB9GEJ5AA8R7K867KB` Genesis detects conflicting decisions across change sets in flight and surfaces them to the owner.
