# conduct-handback-routes

decision · `01M7MJXKKW7XVHWGYFHWAW2QAD`

## reasoning

Routes and their order come from internal/workflow/builtin/ship-change.toml's on.* targets and comments, its conduct handback template (spawn the role, never fix it yourself), the handback.md shape in internal/run/brief.go, and the gate-fix, ci-fix and conflict-fix topics' own scopes; rhiza check and signoff have no dedicated fixer role, so the generic fixer takes them, as the fixer topic fixes findings given in its prompt. A wait-ci timeout means CI never settled (a red CI settles the wait and fails at merge), so ci-fix, which reproduces the tests and findings of rhiza ship check --json, has nothing to act on and the timeout gets no fixer.

## statement

A conductor handed back by ship-change reads handback.md (Failing step, Outcome, Failing commands with their kept output, Output tail, Cause, Failed CI runs, Conflicted files) and spawns, never fixes: gate's rhiza run build-check a verifier for the leaf lacking independent Evidence (a builder retry if it fails); gate's rhiza run review-check a critic for the check with an open blocker or major, then a fixer for confirmed findings and the checks again; gate's rhiza check a fixer with the output as its findings; gate's rhiza ship gate and the push step's rhiza ship open-pr (its pre-push check) the gate-fix role; signoff changes a fixer with the owner's requested changes, then the checks; merge's settled red CI the ci-fix role; wait-ci's timeout spawns no role, since CI never settled and there are no failing tests or findings to fix: the conductor reads rhiza ship check --json, changes nothing and finishes, so the run goes through gate, signoff, push and a fresh CI wait, and a CI that never settles parks at handback_cap; merge's conflicting PR or moved main rhiza ship rebase --no-push (never plain ship rebase), then the conflict-fix role for a code conflict it leaves; every route lands and renders again, then rhiza session done. After handback_cap hand-backs (3 unless set) the next red parks the run stopped with an inbox item, and the steps after the merge never hand back.

## wording

Hand-back routes as ship-change and the fixer role topics give them (derived from the code, autonomous plan run)

## Links

### Out

- derived-from [conduct-guide-selfcontained-plan-2026-11-17](../source/conduct-guide-selfcontained-plan-2026-11-17.md)

### In

- decided-by [conduct-handback-topic](../capability/rhiza/guide/agent-guide/conduct-handback-topic.md)

## History

- `01M7MJN83Y6EF727C2T0AF3Y5D` The rhiza-conduct skill and its guide topics are enough on their own: a conductor outside this repo has no Rhiza source to read, so everything a conductor looked up in the code here (step roles, handback routes, workflow-new defaults, check roles) is in the guide
