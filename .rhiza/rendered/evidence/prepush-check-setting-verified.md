# prepush-check-setting-verified

evidence · `01M770A94PB89SJ5KY4QDRANTK` · actor [claude-code-ff346ed0](../actor/claude-code-ff346ed0.md)

## result

PASS

## statement

Verifier claude-code-ff346ed0 at e292baaae: go test -trimpath ./cmd/rhiza -run 'TestScripts/(ship-prepush-check|describe-golden)$' PASS. With a scratch build in a hermetic repo, rhiza daemon settings --json lists rhiza.ship.prePushCheck between localGate and generate as {scope: project, default: "", value: "", source: default}; after rhiza daemon set rhiza.ship.prePushCheck 'make self-check' it lists value 'make self-check', source 'git config'. docs/reference/config.md and cli.md carry the key; rhiza docs regenerates them with no diff. CLAUDE.md documents setting it to make self-check for this repo. make self-check here fails only on stale .rhiza/rendered (the builder's landed actor record and the verifier's untracked one), which the run's build/render step rewrites; in a clean worktree at the commit, after rhiza render, rhiza check exits 0 with no docs failure. Nit, not blocking: in internal/api/daemon_settings.go the old generate comment now sits above the prePushCheck comment instead of above the generate entry.

## Links

### Out

- verifies [prepush-check-setting](../capability/rhiza/orchestration/shipping/wait-green-pr/prepush-check-setting.md)

## History

- `01M76RH0YQYW1SMCEE5189QDJT` Every run checks the model and generated files locally before it pushes, whatever CI's state, so model and generated-file errors never reach CI. Owner, 2026-11-11, after 'I'm still seeing more ci failures than i would like' and the orchestrator's breakdown, chose option A: "A" (require a pre-push self-check in every run). Evidence: of ~44 PR CI runs on 2026-11-11, three failed on errors a one-minute local check catches: PR #442 lint, stale docs/reference (builders added settings without rhiza docs); agent docs plan PR #456, 'duplicate handle milestone-32'; release plan PR #457, a log entry overwritten by a shell command ('not canonical'). Today ship.gate (internal/api/ship_gate.go) runs the project's rhiza.ship.localGate (this repo: 'make check' in git config) only while CI is unstable, and the plan workflow has no gate at all. Wanted: (1) a fast pre-push check that always runs before any run pushes (plan, plan-milestone, ship-change and its wrappers such as bugfix, every push including fixer and ci-fix pushes and merges of main): rhiza check on the run worktree (model, log canonicity, stale render) plus the project's generated-file command; (2) the generated-file command is a project setting, language neutral (this repo: make self-check, which also rebuilds rhiza and checks docs/reference), so other projects set their own or none; (3) the heavier rhiza.ship.localGate keeps its current meaning (only while CI is unstable); (4) a failed pre-push check goes back to the step's fixer with the failing lines instead of pushing, within the step's cap; (5) the run journal and rhiza run status show the pre-push check's result. Coordinate with milestone 32's language-neutral settings (01M763YX).
