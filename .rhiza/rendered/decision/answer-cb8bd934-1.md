# answer-cb8bd934-1

decision · `01M5Z96GNPBDMZK3YCFVNRP7ME`

## options

```json
[
  "Workflow step kind (Recommended)",
  "Inline copies",
  "Nested runs"
]
```

## picked

Workflow step kind (Recommended)

## question

How should ship-change run Rhiza's build and review in place of the shipit agent step?

## session

cb8bd934-a417-493a-aebe-01a17e247ca8

## statement

How should ship-change run Rhiza's build and review in place of the shipit agent step? — Workflow step kind (Recommended) This departs from four-workflows, which has ship-change's agent step run build then review: ship-change has no agent step, it runs build and review as two `workflow` steps in its own run and journal.

## wording

the agent's option, picked by the owner in session cb8bd934-a417-493a-aebe-01a17e247ca8

## Links

### Out

- relates-to [four-workflows](four-workflows.md) — refines four-workflows: ship-change runs build and review as two workflow steps, not through an agent step
- derived-from [ship-change-plan](../source/ship-change-plan.md)

### In

- decided-by [ship-through-rhiza](../capability/ship-through-rhiza.md)
- decided-by [workflow-step](../capability/ship-through-rhiza/workflow-step.md)
- decided-by [ship-change-file](../capability/ship-through-rhiza/ship-change-file.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
