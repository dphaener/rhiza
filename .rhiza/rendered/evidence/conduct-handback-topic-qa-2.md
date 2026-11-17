# conduct-handback-topic-qa-2

evidence · `01M7MMRHMW3YFNJSPD9K8KMKDR` · actor [claude-code-d305ef13](../actor/claude-code-d305ef13.md)

## result

PASS

## statement

QA cycle 2: with a fresh scratch build, rhiza guide conduct-handback prints handback.md's parts and a route for gate (four commands), signoff, push, wait-ci (timeout, no role) and merge (ci-fix; rhiza ship rebase --no-push then conflict-fix), handback_cap 3 and main-ci/close-out/ship/deploy-check never handing back, all matching internal/workflow/builtin/ship-change.toml; conduct's Read first names it; 629 tokens. In a throwaway worktree: adding a gate command, giving main-ci on.timeout=handback:conduct, adding a main-ci route, and naming ci-fix under wait-ci each fail TestConductHandback naming the step, command or role; dropping it from conduct's Read first fails TestConductReadsHandback. go test ./internal/guide and TestScripts/guide pass.

## Links

### Out

- verifies [conduct-handback-topic](../capability/rhiza/guide/agent-guide/conduct-handback-topic.md)

## History

- `01M7MJN83Y6EF727C2T0AF3Y5D` The rhiza-conduct skill and its guide topics are enough on their own: a conductor outside this repo has no Rhiza source to read, so everything a conductor looked up in the code here (step roles, handback routes, workflow-new defaults, check roles) is in the guide
