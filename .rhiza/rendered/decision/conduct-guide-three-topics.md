# conduct-guide-three-topics

decision · `01M7MJXKHBWGCM2AM6BBC1256A`

## reasoning

conduct is 690 and workflows 883 of the 1000-token TopicCap, so neither can hold the material; the facts live in internal/workflow/spec.go, internal/workflow/builtin/ship-change.toml, internal/run/brief.go, internal/api/run_spawn.go and internal/api/workflow_new.go, and a test over those tables keeps the guide from drifting from them.

## statement

What a conductor looked up in Rhiza's code goes into three new guide topics, conduct-roles (spawn roles and check roles), conduct-handback (hand-back routes) and workflow-keys (every workflow key, step kind and default, and what rhiza workflow new writes), each named in the conduct topic's Read first and each pinned to the code by a guide test that iterates the code's own tables, so a key, role, declared check or route added to the code without the guide fails the test; the conduct and workflows topics are not grown past TopicCap to hold them.

## wording

Three new tested topics the conduct topic points to (derived from TopicCap and the code, autonomous plan run)

## Links

### Out

- derived-from [conduct-guide-selfcontained-plan-2026-11-17](../source/conduct-guide-selfcontained-plan-2026-11-17.md)

### In

- decided-by [conduct-roles-topic](../capability/rhiza/guide/agent-guide/conduct-roles-topic.md)
- decided-by [conduct-handback-topic](../capability/rhiza/guide/agent-guide/conduct-handback-topic.md)
- decided-by [workflow-keys-topic](../capability/rhiza/guide/agent-guide/workflow-keys-topic.md)

## History

- `01M7MJN83Y6EF727C2T0AF3Y5D` The rhiza-conduct skill and its guide topics are enough on their own: a conductor outside this repo has no Rhiza source to read, so everything a conductor looked up in the code here (step roles, handback routes, workflow-new defaults, check roles) is in the guide
