# conduct-handback-topic-qa-1

evidence · `01M7MMMFXK6WVS8KKGY73HWQVZ` · actor [claude-code-22634a90](../actor/claude-code-22634a90.md)

## result

PASS

## statement

QA with a fresh build: rhiza guide conduct-handback prints the topic (629 tokens), named in conduct's Read first; its handback.md parts match internal/run/brief.go's headings; its routes (gate, signoff, push, wait-ci, merge) are exactly ship-change's handback:conduct steps and its four gate commands match the gate step; rhiza ship rebase --no-push exists per its --help; handback_cap 3 matches ship-change; TestConductHandback* and TestScripts/guide pass

## Links

### Out

- verifies [conduct-handback-topic](../capability/rhiza/guide/agent-guide/conduct-handback-topic.md)

## History

- `01M7MJN83Y6EF727C2T0AF3Y5D` The rhiza-conduct skill and its guide topics are enough on their own: a conductor outside this repo has no Rhiza source to read, so everything a conductor looked up in the code here (step roles, handback routes, workflow-new defaults, check roles) is in the guide
