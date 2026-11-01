# roadmap-app-new-work

capability · `01M686V49K74FDKG4YV2BDJCSN` · parent [roadmap-app-start-runs](../roadmap-app-start-runs.md)

## acceptance

A Go test with a fake Starter.Exec: m, typed text, enter, y gives one rex window new running run start --workflow plan --new <intent> --milestone; n gives the same without --milestone; esc or n at the prompt and empty text call nothing

## statement

In the roadmap app, m and n prompt for a one-line intent, ask y/n, then open a Rex tab running a plan run for a new milestone or a new change set, with no change set selected

## Links

### Out

- decided-by [new-work-through-run-start](../../../decision/new-work-through-run-start.md)

### In

- verifies [roadmap-app-new-work-verified](../../../evidence/roadmap-app-new-work-verified.md)
- verifies [roadmap-app-new-work-qa-1](../../../evidence/roadmap-app-new-work-qa-1.md)
- verifies [pr-148-new-work-from-roadmap](../../../evidence/pr-148-new-work-from-roadmap.md)

## History

- `01M686RVAD1EHGERBBDJN3S45G` The roadmap app starts a plan run for a new milestone or a new change set from a typed one-line intent
