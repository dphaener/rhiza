# workflow-step-limits-plan

source · `01M6504MNQD634VWNBCG720D55`

## quotes

```json
[]
```

## raw_sha256

59f4c8c1083d2883235d83a31457d256688e594fb106b48b5517e8b5e74db6b7

## session

d17e5b88-5c8f-4163-b32f-01a18a012af1

## summary

Autonomous plan run (no owner present) for change set 01M60AW, intent: record the workflow step's limits as constraints. Two limits are recorded, each as it stands in the code: a workflow runs another workflow at most once (internal/workflow/include.go), and a workflow step takes no prompt, agent, permission_mode, run, until, until_when, timeout, interval or max_visits keys, and needs workflow = a workflow's name (buildWorkflow in internal/workflow/workflow.go). No owner words exist for this plan, so no quotes are recorded. Open: none recorded here.

## title

Workflow step limits, planned without an owner

## Links

### In

- derived-from [workflow-runs-another-once](../constraint/workflow-runs-another-once.md)
- derived-from [workflow-step-no-session-keys](../constraint/workflow-step-no-session-keys.md)

## History

- `01M60AWXJ785WZCDW1G1KA63BR` Record the workflow step's limits as constraints: a workflow runs another at most once, and a workflow step takes no session, condition or max_visits keys
