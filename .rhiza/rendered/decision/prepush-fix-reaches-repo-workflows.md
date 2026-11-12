# prepush-fix-reaches-repo-workflows

decision · `01M76SKG2THQYFMKYAYC9GVC7X`

## reasoning

The finding is right on the mechanics: workflow-overrides makes a repo file shadow the built-in for every bare-name use (including bugfix's workflow step), and include.go rejects a repo ship-change including ship-change as a self-include, so 01M763YX's repo workflows (rhiza-own-workflows-go-cleanup, ship-head-build-step) must be full copies. Then routing added only in internal/workflow/builtin never reaches Rhiza's own runs, which are the runs whose CI failures motivated this change set (intent points 1 and 4); the push itself would still be refused by shipPush, but the run would end instead of reaching a fixer. Options weighed: (a) wait for or depend on 01M763YX: couples two independent change sets and does not fix the case where 01M763YX lands first; (b) add a builtin: include prefix to the resolver: new language surface outside this change set's intent and scope; (c) port the routing into whatever repo copies exist and add a self-check guard: smallest, works in either landing order, and the guard turns the silent gap rhiza overlaps cannot see (different nodes) into a failing make self-check for whichever change set lands second. The intent asks to coordinate with milestone-32; this makes that coordination concrete without changing 01M763YX's nodes.

## statement

The pre-push fix routing must reach Rhiza's own repo-level workflow copies, not only the built-ins. Since a repo .rhiza/workflows/<name>.toml shadows the built-in and cannot include it, change set 01M76RH0 owns keeping them in step: when ship-change-prepush-fix and plan-prepush-fix are built, any .rhiza/workflows/ship-change.toml, plan.toml or plan-milestone.toml present on main at that time gets the same on.fail -> ship gate --failed triage -> gate-fix / gate-fix-merge / prepush-fix routing and steps as the built-in it copies; and both leaves add a check, run by make self-check (via rhiza workflow check or an internal/workflow test over the repo's .rhiza/workflows), that fails if a repo copy of ship-change, plan or plan-milestone lacks that routing, so whichever of 01M76RH0 and 01M763YX lands second is caught and must carry the routing into the copies. Both leaves' acceptance is extended to cover this. 01M76RH0 does not wait on 01M763YX and does not add a builtin: include prefix.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m76rgz9k7zmdyec31e5kq7ta-r3-f1](../question/review-01m76rgz9k7zmdyec31e5kq7ta-r3-f1.md)

## History

- `01M76RH0YQYW1SMCEE5189QDJT` Every run checks the model and generated files locally before it pushes, whatever CI's state, so model and generated-file errors never reach CI. Owner, 2026-11-11, after 'I'm still seeing more ci failures than i would like' and the orchestrator's breakdown, chose option A: "A" (require a pre-push self-check in every run). Evidence: of ~44 PR CI runs on 2026-11-11, three failed on errors a one-minute local check catches: PR #442 lint, stale docs/reference (builders added settings without rhiza docs); agent docs plan PR #456, 'duplicate handle milestone-32'; release plan PR #457, a log entry overwritten by a shell command ('not canonical'). Today ship.gate (internal/api/ship_gate.go) runs the project's rhiza.ship.localGate (this repo: 'make check' in git config) only while CI is unstable, and the plan workflow has no gate at all. Wanted: (1) a fast pre-push check that always runs before any run pushes (plan, plan-milestone, ship-change and its wrappers such as bugfix, every push including fixer and ci-fix pushes and merges of main): rhiza check on the run worktree (model, log canonicity, stale render) plus the project's generated-file command; (2) the generated-file command is a project setting, language neutral (this repo: make self-check, which also rebuilds rhiza and checks docs/reference), so other projects set their own or none; (3) the heavier rhiza.ship.localGate keeps its current meaning (only while CI is unstable); (4) a failed pre-push check goes back to the step's fixer with the failing lines instead of pushing, within the step's cap; (5) the run journal and rhiza run status show the pre-push check's result. Coordinate with milestone 32's language-neutral settings (01M763YX).
