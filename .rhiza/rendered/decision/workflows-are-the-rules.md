# workflows-are-the-rules

decision · `01M6VY483G51DG4WY4PR58AMXV`

## options

```json
[
  "rules as workflow policy with built-in defaults",
  "rules hard-coded in the daemon"
]
```

## quote

"everything should run through a workflow. those are the rules that are defined by the user/owner. rhiza ships with sensible defaults but should be flexible. the only thing we enforce is the references between nodes." (owner, 2026-11-07); also "we should be able to have bug fixes without a milestone. smaller features without a milestone. hell even no milestones at all if people don't want that. autonomy should be configurable on a workflow by workflow basis." (owner, 2026-11-07). Agent's reading, confirmed by the owner ("yes this sounds correct"): the append-only log and change-set lifecycle stay as mechanism; who may accept, the daemon's milestone skip, autonomy and the milestone-needs-a-human rule become defaults in built-in workflows.

## statement

Everything Rhiza does runs through a workflow: the rules (which workflow a request or change set gets, whether work without a milestone is built and in what order, each stage's autonomy, whether milestone planning needs the owner) are workflow and repo policy the owner defines. Rhiza ships them as built-in defaults; the only rule it enforces is that references between nodes hold.

## Links

### In

- decided-by [workflows](../capability/rhiza/workflows.md)

## History

- `01M6VY4688KRW7ZMSTA847M80V` Record the owner's principle: every rule is workflow policy; Rhiza enforces only references between nodes
