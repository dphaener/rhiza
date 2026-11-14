# plan-merge-race-research

source · `01M7AQEFZDTZ6MYYN5RV0W9TFQ`

## quotes

```json
[]
```

## raw_sha256

b7703d2e12b91027bf8ef0dbc9f24509c02037883e3d2f1eb2e3de2814a3a77a

## session

f285a605-5558-4e2d-80e5-01a1d5746e62

## summary

Autonomous plan run, no owner. Run 01M7AH6RHKMZRAHGJKJFM4VTNT (built-in plan workflow, PR #551) reached merge three times between 14:34 and 14:44 on 2026-11-13; each time ship merge refused with 'origin/main has moved past PR #551's checked head' (ecd6277, 4931719, 6ce8f26) because other plan runs merged their render PRs while the rebased head waited 2-7 minutes for CI. The fourth rebase went over the rebase step's max_visits 3 and the watchdog filed this bug with the PR green and unmerged. The plan and plan-milestone workflows cap rebase, wait-rebased and merge at 3 visits and rebased-green at 4. Decision plan-merge-loops-to-rebase set those caps to 'ship-change's'; ship-change's rebase loop (rebase, wait-rebased, rebased-green, rebase-merge, wait-rebased-merge, rebased-merge-green, merge) has since been raised to 6 visits each, and the plan workflows never followed. No in-flight change set overlaps.

## title

Plan run lost the merge race three times (autonomous plan-bug research)

## Links

### In

- derived-from [plan-merge-loop-cap-six](../decision/plan-merge-loop-cap-six.md)

## History

- `01M7AQ89397MPNWGFSAYG3BK77` Fix watchdog unmerged PR at step rebase: run 01M7AH6RHKMZRAHGJKJFM4VTNT; log: PR green at 6ce8f26 and unmerged after step rebase failed: step "rebase" went over max_visits 3
