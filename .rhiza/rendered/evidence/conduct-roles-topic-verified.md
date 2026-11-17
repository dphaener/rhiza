# conduct-roles-topic-verified

evidence · `01M7MKQQ5E00ZRZV9P85CM129R` · actor [claude-code-29e98897](../actor/claude-code-29e98897.md)

## result

PASS

## statement

At 97ee44c3b re-ran the Checks block: go test -trimpath ./internal/guide/... -run 'TestConductRoles|TestConduct|TestList' -v (TestConductRoles, TestConductReadsRoles, TestConductRolesCatchesDrift, TestConduct, TestListOrder, TestList PASS), go test ./internal/guide/... ok, go vet ./internal/guide ok, go test ./internal/api -run TestRunSpawnCheck PASS, go test ./cmd/rhiza -run TestScripts/guide ok. Diff read: conduct-roles.md lists the ten roles each as rhiza guide <role>, the --brief/--leaf/--check/--model flags, [[step.<name>.check]] conductor steps and the arbiter, RHIZA_CHECK/RHIZA_CYCLE, the refusal text matching run_spawn.go ('name one of its checks: plan, code, qa'), and ship-change's plan/code/qa vs the classic compliance/quality/model; conduct.md Read first names rhiza guide conduct-roles and stays under TopicCap; the test parses built-in ship-change via the workflow package and requires check-<name> per declared check, and TestConductRolesCatchesDrift covers the failure cases (unknown role check-scope, undeclared lint check, missing check-qa). No spawn, workflow or built-in workflow changes; no TODOs.

## Links

### Out

- verifies [conduct-roles-topic](../capability/rhiza/guide/agent-guide/conduct-roles-topic.md)

## History

- `01M7MJN83Y6EF727C2T0AF3Y5D` The rhiza-conduct skill and its guide topics are enough on their own: a conductor outside this repo has no Rhiza source to read, so everything a conductor looked up in the code here (step roles, handback routes, workflow-new defaults, check roles) is in the guide
