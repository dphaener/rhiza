# workflow-step-no-session-keys

constraint · `01M6504V7GEH5DQDQEC4200VRS`

## acceptance

A test refuses each of those keys on a workflow step, with a message naming the key

## statement

A workflow step takes no session, condition or max_visits keys: prompt, agent, permission_mode, run, until, until_when, timeout, interval and max_visits are refused, and it needs workflow = the name of a workflow

## Links

### Out

- constrains [workflow-definitions](../capability/rhiza/workflows/workflow-definitions.md)
- derived-from [workflow-step-limits-plan](../source/workflow-step-limits-plan.md)

## History

- `01M60AWXJ785WZCDW1G1KA63BR` Record the workflow step's limits as constraints: a workflow runs another at most once, and a workflow step takes no session, condition or max_visits keys
