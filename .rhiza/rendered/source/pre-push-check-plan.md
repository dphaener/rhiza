# pre-push-check-plan

source · `01M76RSBKCAF91E8NY2JS692TY`

## quotes

```json
[]
```

## raw_sha256

999f774c8cce778291e230a93a4cbd2a14e45d723bff5891636e86f9af3a4820

## session

aa618fe3-a94f-4379-98c2-01a1cd88ca37

## summary

Autonomous plan run for change set 01M76RH0, no owner. The planner read the intent (owner chose option A on 2026-11-11: a pre-push self-check in every run), the push paths in code (shipPush in internal/api/ship_close.go used by ship open-pr and the land push; ship rebase's push in ship_rebase.go), ship.gate and rhiza.ship.localGate, ship.generate, and the plan, plan-milestone and ship-change workflows. It raised one batched question (which intent, where to enforce, the setting's name, where a plan's failure goes). The arbiter answered with Decision pre-push-check-shape: serve milestone-25 like ship-local-gate did; enforce in the shared push path, running rhiza check then the project's command, refusing the push on failure; setting rhiza.ship.prePushCheck, empty default, make self-check for this repo; localGate keeps its meaning; ship-change routes a failed pre-push check to gate-fix within its cap; plan and plan-milestone get a prepush-fix agent step (max_visits 2) back to open-pr. The planner then wrote five leaves under wait-green-pr.

## title

Plan of the pre-push check in every run (autonomous, arbiter)

## Links

### In

- derived-from [pre-push-check-shape](../decision/pre-push-check-shape.md)

## History

- `01M76RH0YQYW1SMCEE5189QDJT` Every run checks the model and generated files locally before it pushes, whatever CI's state, so model and generated-file errors never reach CI. Owner, 2026-11-11, after 'I'm still seeing more ci failures than i would like' and the orchestrator's breakdown, chose option A: "A" (require a pre-push self-check in every run). Evidence: of ~44 PR CI runs on 2026-11-11, three failed on errors a one-minute local check catches: PR #442 lint, stale docs/reference (builders added settings without rhiza docs); agent docs plan PR #456, 'duplicate handle milestone-32'; release plan PR #457, a log entry overwritten by a shell command ('not canonical'). Today ship.gate (internal/api/ship_gate.go) runs the project's rhiza.ship.localGate (this repo: 'make check' in git config) only while CI is unstable, and the plan workflow has no gate at all. Wanted: (1) a fast pre-push check that always runs before any run pushes (plan, plan-milestone, ship-change and its wrappers such as bugfix, every push including fixer and ci-fix pushes and merges of main): rhiza check on the run worktree (model, log canonicity, stale render) plus the project's generated-file command; (2) the generated-file command is a project setting, language neutral (this repo: make self-check, which also rebuilds rhiza and checks docs/reference), so other projects set their own or none; (3) the heavier rhiza.ship.localGate keeps its current meaning (only while CI is unstable); (4) a failed pre-push check goes back to the step's fixer with the failing lines instead of pushing, within the step's cap; (5) the run journal and rhiza run status show the pre-push check's result. Coordinate with milestone 32's language-neutral settings (01M763YX).
