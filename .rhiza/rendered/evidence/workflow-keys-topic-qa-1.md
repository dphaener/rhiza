# workflow-keys-topic-qa-1

evidence · `01M7MMMG05MQ8JWDHXKJMHS0A3` · actor [claude-code-22634a90](../actor/claude-code-22634a90.md)

## result

PASS

## statement

QA with a fresh build in a scratch repo: rhiza guide workflow-keys prints the topic (910 tokens), named by conduct and workflows (which lists gate); rhiza workflow new demo writes start=build, agent step build, command step test with a failing placeholder, max_visits=3, on.fail=handback:build as the topic says; a second new demo is refused (exit 3, file left); new --from ship-change copies with name set; --from nosuch fails cleanly (exit 2); workflow check demo validates with the placeholder warning; TestWorkflowKeys*, TestWorkflowsText and TestScripts/guide pass

## Links

### Out

- verifies [workflow-keys-topic](../capability/rhiza/guide/agent-guide/workflow-keys-topic.md)

## History

- `01M7MJN83Y6EF727C2T0AF3Y5D` The rhiza-conduct skill and its guide topics are enough on their own: a conductor outside this repo has no Rhiza source to read, so everything a conductor looked up in the code here (step roles, handback routes, workflow-new defaults, check roles) is in the guide
