# workflow-runs-another-once

constraint · `01M6504ST566B8PMBWN10XJEAV`

## acceptance

A test refuses a definition whose two steps both run the same workflow, naming both steps

## statement

A workflow runs another workflow at most once: two workflow steps of one workflow may not name the same workflow

## Links

### Out

- constrains [workflow-definitions](../capability/rhiza/workflows/workflow-definitions.md)
- derived-from [workflow-step-limits-plan](../source/workflow-step-limits-plan.md)

## History

- `01M60AWXJ785WZCDW1G1KA63BR` Record the workflow step's limits as constraints: a workflow runs another at most once, and a workflow step takes no session, condition or max_visits keys
