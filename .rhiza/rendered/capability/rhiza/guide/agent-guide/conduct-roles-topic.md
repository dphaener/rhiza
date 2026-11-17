# conduct-roles-topic

capability · `01M7MJXKRW5GHYEGJPFV4GY7XT` · parent [agent-guide](../agent-guide.md)

## acceptance

go test -trimpath ./internal/guide -run 'TestConductRoles|TestConduct|TestList' shows the conduct-roles topic listed and under TopicCap, the conduct topic's Read first naming rhiza guide conduct-roles and staying under TopicCap, and conduct-roles naming each of the ten roles as rhiza guide <role> with every one a topic Lookup finds, plus --brief, --leaf, --check, --model, RHIZA_CHECK, RHIZA_CYCLE, [[step.<name>.check]], the arbiter and the refusal of an undeclared --check; the test parses the built-in ship-change with the workflow package and requires check-<name> for every check its conduct step declares; (failure case) the test fails naming the role when the topic lists a role with no guide topic of that name, or a check ship-change declares has no check-<name> line; go test -trimpath ./internal/api -run TestRunSpawnCheck still passes (an undeclared --check refused naming plan, code, qa); go test -trimpath ./cmd/rhiza -run TestScripts/guide passes.

## files

internal/guide/text/conduct-roles.md, internal/guide/conduct_roles_test.go, internal/guide/guide.go, internal/guide/guide_test.go, internal/guide/text/conduct.md, internal/guide/conduct_test.go, cmd/rhiza/testdata/script/guide.txtar

## must_not

Change rhiza run spawn, the workflow package or any built-in workflow: guide text and its tests only; grow the conduct topic past TopicCap; change what the review index topic says

## statement

The guide has a conduct-roles topic, named in the conduct topic's Read first, listing every role a conductor spawns with rhiza run spawn (builder, verifier, check-plan, check-code, check-qa, critic, fixer, gate-fix, ci-fix, conflict-fix), each with the rhiza guide topic and skill it follows (its own name), when the conductor spawns it and the flags it takes (--brief always, --leaf for builder and verifier, --check <check> for a critic, --model optional), and the check roles: an agent step that declares [[step.<name>.check]] tables is a conductor step whose session and children put questions to the arbiter; a check-<check> role sets the child's RHIZA_CHECK and RHIZA_CYCLE; --check naming a check the step does not declare is refused with the declared names; the conductor default ship-change declares plan, code and qa, while compliance, quality and model belong to the classic workflows.

## Links

### Out

- decided-by [conduct-guide-three-topics](../../../../decision/conduct-guide-three-topics.md)

### In

- verifies [conduct-roles-topic-verified](../../../../evidence/conduct-roles-topic-verified.md)
- verifies [conduct-roles-topic-qa-1](../../../../evidence/conduct-roles-topic-qa-1.md)
- verifies [conduct-roles-topic-qa-2](../../../../evidence/conduct-roles-topic-qa-2.md)

## History

- `01M7MJN83Y6EF727C2T0AF3Y5D` The rhiza-conduct skill and its guide topics are enough on their own: a conductor outside this repo has no Rhiza source to read, so everything a conductor looked up in the code here (step roles, handback routes, workflow-new defaults, check roles) is in the guide
