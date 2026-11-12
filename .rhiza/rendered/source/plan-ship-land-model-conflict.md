# plan-ship-land-model-conflict

source · `01M77N9K9BS1PDMZD2GDFDAT7K`

## quotes

```json
[]
```

## raw_sha256

aca44ca9cc9a9c52a330adda6b80b296108df831d1593935648fc325aec50052

## session

a86fdccf-e6b8-4c91-ad56-01a1cf50f590

## summary

Autonomous plan run, no owner. Diagnosis: ship-change run 01M7786QKAZMH7PBJZP7QQA6HX failed at step land because change set 01M763YX47S5XY0W3P8FVJ54S8 (Rhiza's Go behaviour moves, PR 475) updated root-daemon-settings-scope's statement on a base that change set 01M763Z16T0ZGQMY1EC2Y81V0K (auto-update, PR 473) had since superseded with its own statement write. rhiza land refused it as a model conflict (exit 3). land's on.fail goes to land-triage, which runs rhiza ship gate --failed and passes only a recorded pre-push failure, so the run ended and the watchdog filed this bug. The orchestrator repaired it by hand with rhiza set --rebase, keeping both edits (commit 480b6acf3), and the land retry passed. The planner asked how ship-change should handle such a refusal: an agent repair step, an inbox question, or code auto-rebase. The arbiter decided ship-land-model-conflict-rebase: land-triage tells a model conflict apart from other failures and routes it to a model-conflict fixer agent that rebases each stale update with rhiza set --rebase so both edits hold, commits with Change: and Actor: trailers, and returns through gate-land to land. The fixer has a max_visits cap and asks one blocking question on a real clash. Every other land failure ends the run as now. The planner split the work into two leaves: ship land records the conflict for a triage check (internal/api), and the built-in workflow's routes and fixer step.

## title

Plan: ship-change recovers from a land refused for a model conflict (autonomous)

## Links

### In

- derived-from [ship-land-model-conflict-rebase](../decision/ship-land-model-conflict-rebase.md)

## History

- `01M77N1GAGBPBP7EH4EV2CT023` Fix watchdog unmerged PR at step land: run 01M7786QKAZMH7PBJZP7QQA6HX; log: R head's build: conflict: cannot land change set 01M763YX47S5XY0W3P8FVJ54S8: model: conflict: base 01M6RN85C17S5VE1A6SQVD7WE1 of 01M6RN85C17S5VE1A6SNZAXK89 is stale: op 01M77C4T7EBWGBK547B9S31A74 of change set 01M763Z16T0ZGQMY1EC2Y81V0K since wrote field:statement: /home/rhiza/code/rhiza/.git/rhiza/ship/01M763YX47S5XY0W3P8FVJ54S8/rhiza land --change 01M763YX47S5XY0W3P8FVJ54S8 --json: exit status 3
