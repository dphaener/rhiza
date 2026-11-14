# handback-split-route-deliver

decision · `01M7AHW75ZQHSCV0CGQ3G4CGD9`

## statement

The hand-back primitive is two leaves: routing (the handback:<agent step> target, the per-run count, handback_cap and parking) and delivery (handback.md, the handback text, messaging a live session or resuming an ended one); rhiza workflow check's warning about an uncapped hand-back loop stays with the draft change set that promises it.

## wording

derived from the model by the planner of an autonomous plan run

## Links

### Out

- derived-from [plan-daemon-offers-spawn](../source/plan-daemon-offers-spawn.md)

### In

- decided-by [handback-route](../capability/rhiza/workflows/workflow-runs/handback-route.md)
- decided-by [run-handback](../capability/rhiza/workflows/workflow-runs/run-handback.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
