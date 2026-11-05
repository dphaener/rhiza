# plan-merge-rebases-interview

source · `01M6MT96GM8XF44X2PX1WS17Y1`

## quotes

```json
[]
```

## raw_sha256

5735dda0b4cdec9aa5f1221631fee2ab62bb7e6e92113fcb1b08debec3c73ad6

## session

912154e0-797e-473f-8c79-01a1a9a34415

## summary

Autonomous plan run, no owner present. The change set's intent names the fix: the plan workflow's merge step, which fails when origin/main moved past the checked head, should loop back to rebase within a cap, as ship-change's merge does (on.fail = rebase-merge, max_visits = 3). The planner placed the change set under milestone 25 (empty backlog: bugs and polish), as the sibling fix 01M6M1M1 is, and planned one leaf under plan-workflow: plan.toml's merge step goes back to rebase on fail, the rebase/wait-rebased/rebased-green/merge loop capped at 3 visits like ship-change, proven by a Go workflow transition test and a testscript modelled on run-ship-merge-main-moved.txtar. Nothing left open.

## title

Plan of Fix Plan runs fail at merge (autonomous)

## Links

### In

- derived-from [plan-merge-loops-to-rebase](../decision/plan-merge-loops-to-rebase.md)

## History

- `01M6MT67M8VH509KK24BP65QRG` Fix Plan runs fail at merge when main moves: the plan workflow's merge step ends the run with 'Not merged: origin/main has moved past PR's checked head; rebase and check again' (run 01M6MM3J, PR #272, 2026-11-05) instead of going back to rebase and wait for green as ship-change does, so any merge landing between rebased-green and merge strands a plan PR that then also conflicts on .rhiza/rendered. The merge step should loop back to rebase, within a cap.
