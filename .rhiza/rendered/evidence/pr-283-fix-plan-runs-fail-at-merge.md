# pr-283-fix-plan-runs-fail-at-merge

evidence · `01M6MX621P9B2RKE48FN6CEBYV` · actor [dphaener](../actor/dphaener.md)

## head

64afccee2865ca6cae31fcf72caf1f90919ad0d8

## pr

https://github.com/dphaener/rhiza/pull/283

## result

1 check passed: check SUCCESS

## statement

Pull request #283 (https://github.com/dphaener/rhiza/pull/283) passed CI at 64afcce

## status

verified

## Links

### Out

- verifies [plan-run-merge-rebases](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-run-merge-rebases.md)

## History

- `01M6MT67M8VH509KK24BP65QRG` Fix Plan runs fail at merge when main moves: the plan workflow's merge step ends the run with 'Not merged: origin/main has moved past PR's checked head; rebase and check again' (run 01M6MM3J, PR #272, 2026-11-05) instead of going back to rebase and wait for green as ship-change does, so any merge landing between rebased-green and merge strands a plan PR that then also conflicts on .rhiza/rendered. The merge step should loop back to rebase, within a cap.
