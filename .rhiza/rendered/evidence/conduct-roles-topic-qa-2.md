# conduct-roles-topic-qa-2

evidence · `01M7MMRHGH01HKVXRDTZ7A9BGC` · actor [claude-code-d305ef13](../actor/claude-code-d305ef13.md)

## result

PASS

## statement

QA cycle 2: built rhiza fresh to a scratch path (go build -trimpath) and ran it: rhiza guide conduct-roles prints the ten roles with --brief/--leaf/--check/--model, RHIZA_CHECK/RHIZA_CYCLE, [[step.<name>.check]], the arbiter and the refusal text 'name one of its checks: plan, code, qa', matching internal/api/run_spawn.go and rhiza run spawn --help; rhiza guide <role> resolves for all ten roles; rhiza guide conduct's Read first names conduct-roles; guide ls lists it at 708 tokens (40 texts); unknown topics (conduct-role, CONDUCT-ROLES) exit 4 with 'no guide text'. In a throwaway worktree, renaming ci-fix to ci-fixer and declaring a security check in ship-change each fail TestConductRoles naming the role/check. go test ./internal/guide, TestScripts/guide and internal/api TestRunSpawnCheck pass.

## Links

### Out

- verifies [conduct-roles-topic](../capability/rhiza/guide/agent-guide/conduct-roles-topic.md)

## History

- `01M7MJN83Y6EF727C2T0AF3Y5D` The rhiza-conduct skill and its guide topics are enough on their own: a conductor outside this repo has no Rhiza source to read, so everything a conductor looked up in the code here (step roles, handback routes, workflow-new defaults, check roles) is in the guide
