# prepush-check-setting

capability · `01M76RSZVP98R5EET9VWCRT1XC` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

go test -trimpath ./cmd/rhiza -run TestScripts/ship-prepush-check passes, showing rhiza daemon settings --json listing rhiza.ship.prePushCheck as a project setting with an empty default (source default), and rhiza daemon set rhiza.ship.prePushCheck 'make self-check' then listing it with value make self-check from git config; make self-check passes with docs/reference regenerated

## statement

rhiza.ship.prePushCheck is a project setting holding a language-neutral shell command the pre-push check runs after rhiza check; its default is empty (rhiza check alone), it is listed beside rhiza.ship.localGate and rhiza.ship.generate in rhiza daemon settings and docs/reference, and Rhiza's own repo documents setting it to make self-check

## Links

### Out

- decided-by [pre-push-check-shape](../../../../../decision/pre-push-check-shape.md)

### In

- depends-on [prepush-check-before-push](prepush-check-before-push.md)
- verifies [prepush-check-setting-verified](../../../../../evidence/prepush-check-setting-verified.md)
- verifies [prepush-check-setting-qa-1](../../../../../evidence/prepush-check-setting-qa-1.md)
- verifies [pr-469-every-run-checks-the-model-and](../../../../../evidence/pr-469-every-run-checks-the-model-and.md)

## History

- `01M76RH0YQYW1SMCEE5189QDJT` Every run checks the model and generated files locally before it pushes, whatever CI's state, so model and generated-file errors never reach CI. Owner, 2026-11-11, after 'I'm still seeing more ci failures than i would like' and the orchestrator's breakdown, chose option A: "A" (require a pre-push self-check in every run). Evidence: of ~44 PR CI runs on 2026-11-11, three failed on errors a one-minute local check catches: PR #442 lint, stale docs/reference (builders added settings without rhiza docs); agent docs plan PR #456, 'duplicate handle milestone-32'; release plan PR #457, a log entry overwritten by a shell command ('not canonical'). Today ship.gate (internal/api/ship_gate.go) runs the project's rhiza.ship.localGate (this repo: 'make check' in git config) only while CI is unstable, and the plan workflow has no gate at all. Wanted: (1) a fast pre-push check that always runs before any run pushes (plan, plan-milestone, ship-change and its wrappers such as bugfix, every push including fixer and ci-fix pushes and merges of main): rhiza check on the run worktree (model, log canonicity, stale render) plus the project's generated-file command; (2) the generated-file command is a project setting, language neutral (this repo: make self-check, which also rebuilds rhiza and checks docs/reference), so other projects set their own or none; (3) the heavier rhiza.ship.localGate keeps its current meaning (only while CI is unstable); (4) a failed pre-push check goes back to the step's fixer with the failing lines instead of pushing, within the step's cap; (5) the run journal and rhiza run status show the pre-push check's result. Coordinate with milestone 32's language-neutral settings (01M763YX).
