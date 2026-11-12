# ship-land-model-conflict-two-designs

question · `01M77TAX7V4KNR1XANH6BAG8FS`

## statement

Two accepted change sets plan contradictory fixes for a land refused on a model conflict: this one (01M77N83, Decision ship-land-model-conflict-repair: land reports the conflict, land-triage routes to a model-conflict fixer that runs `rhiza set --rebase`, and the gate and rhiza ship check --model-conflict fail on the shipping change set's own conflicts) and 01M77N1G (Decision ship-land-model-conflict-rebase, leaves ship-land-records-model-conflict and ship-change-model-conflict-fix-step: steps land-conflict-triage and model-conflict-fix -> gate-land, a Model-Conflict: trailer, a 3-fix count in ship check --model-conflict, land max_visits 4, gate-land 9). Only one design can ship; the other's capabilities would no longer match the code. Ship-land-model-conflict-repair is meant to supersede ship-land-model-conflict-rebase, a link the model cannot record while 01M77N1G is unlanded. Should 01M77N1G be retired (withdrawn) in favour of this change set, or this one in favour of it?

## Links

### Out

- relates-to [ship-land-model-conflict-repair](../decision/ship-land-model-conflict-repair.md) — The question asks which of two contradictory designs for a land refused on a model conflict ships; this Decision is one of them.

### In

- answers [merge-keeps-main-model-conflict-design](../decision/merge-keeps-main-model-conflict-design.md) — 01M77N1G's design ships: it landed via #484, and this change set's competing implementation is dropped in the merge.

## History

- `01M77N83F027XBGGZPK9SV6YN2` Fix watchdog unmerged PR at step land-triage: run 01M77DWNRDSFY0PEP0S1P58P5W; log: PR green at 52ace47 and unmerged after step land-triage failed: step "land-triage" ended with outcome "fail" | exit 1 | No head the local gate failed on after a rebase is on the branch.
