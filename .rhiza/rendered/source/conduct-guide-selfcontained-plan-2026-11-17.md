# conduct-guide-selfcontained-plan-2026-11-17

source · `01M7MJVAGET23A9G6BHYKKSPYY`

## quotes

```json
[]
```

## raw_sha256

2ad91d75de9ecb9199830b1a7fab1450abe531d06d74b460daf0386ece67127d

## session

d5023e1f-398b-47ad-ba1b-01a1e92ae121

## summary

Autonomous plan run with no owner, so no owner words. Research pass: the conduct topic (690 of its 1000-token cap) names the spawn roles only in passing and says nothing of hand-backs; the workflows topic (883 tokens) lacks the gate step kind, handback:<agent step> routes, handback_cap, stall, cleanup, env, autonomy, an agent step's [[step.<name>.check]] declaration that makes it a conductor step, and the defaults of unset keys and of rhiza workflow new. The facts live only in code: internal/workflow/spec.go (keys and defaults), internal/workflow/builtin/ship-change.toml (routes and declared checks), internal/run/brief.go (handback.md shape), internal/api/run_spawn.go (role to rhiza guide <role>, check-<check> sets RHIZA_CHECK, --check validated against the step's declared checks) and internal/api/workflow_new.go (the skeleton). Neither topic has room for them under TopicCap, so the plan adds three topics the conduct topic points to: conduct-roles, conduct-handback and workflow-keys, each pinned to the code by a test so it cannot drift. The sibling conductor fix 01M7MD2K2BE9 serves m36-conductor, so this change set serves it too. Nothing needed the arbiter.

## title

Planning the self-contained rhiza-conduct guide, autonomous plan run

## Links

### In

- derived-from [conduct-guide-three-topics](../decision/conduct-guide-three-topics.md)
- derived-from [conduct-handback-routes](../decision/conduct-handback-routes.md)
- derived-from [conduct-docs-reply-header-drafted](../decision/conduct-docs-reply-header-drafted.md)

## History

- `01M7MJN83Y6EF727C2T0AF3Y5D` The rhiza-conduct skill and its guide topics are enough on their own: a conductor outside this repo has no Rhiza source to read, so everything a conductor looked up in the code here (step roles, handback routes, workflow-new defaults, check roles) is in the guide
