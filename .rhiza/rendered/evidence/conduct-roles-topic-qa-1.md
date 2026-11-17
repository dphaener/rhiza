# conduct-roles-topic-qa-1

evidence · `01M7MMMFT9WVV1M6H2BJD0DVSY` · actor [claude-code-22634a90](../actor/claude-code-22634a90.md)

## result

PASS

## statement

QA with a fresh build (go build -trimpath to a scratch path): rhiza guide conduct-roles prints the topic (708 tokens, listed by guide ls); rhiza guide conduct's Read first names conduct-roles; all ten roles it lists (builder..conflict-fix) resolve with rhiza guide <role>; its spawn flags match rhiza run spawn --help; the refusal text 'name one of its checks' is in internal/api/run_spawn.go; ship-change declares plan, code, qa; go test ./internal/guide (TestConductRoles*, TestConduct, TestList), ./internal/api -run TestRunSpawnCheck and ./cmd/rhiza TestScripts/guide pass

## Links

### Out

- verifies [conduct-roles-topic](../capability/rhiza/guide/agent-guide/conduct-roles-topic.md)

## History

- `01M7MJN83Y6EF727C2T0AF3Y5D` The rhiza-conduct skill and its guide topics are enough on their own: a conductor outside this repo has no Rhiza source to read, so everything a conductor looked up in the code here (step roles, handback routes, workflow-new defaults, check roles) is in the guide
