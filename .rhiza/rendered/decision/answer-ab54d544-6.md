# answer-ab54d544-6

decision · `01M5YXDQKBNNMN7YAEVERBKTPE`

## options

```json
[
  "Fix commit, same trailer (Recommended)",
  "Autosquash at the end"
]
```

## picked

Fix commit, same trailer (Recommended)

## question

How does a retry commit its fix?

## session

ab54d544-8c72-4c39-ab01-01a17dc7f611

## statement

How does a retry commit its fix? — Fix commit, same trailer (Recommended)

## wording

the agent's option, picked by the owner in session ab54d544-8c72-4c39-ab01-01a17dc7f611

## Links

### Out

- relates-to [build-workflow-shape](build-workflow-shape.md) — Amends build-workflow-shape's "one commit per leaf": the owner picked a fix commit per retry with the same Leaf: trailer over autosquashing to one commit at the end, so a leaf's work is every commit with its Leaf: trailer.
- derived-from [build-workflow-plan](../source/build-workflow-plan.md)

### In

- decided-by [build-leaf-commit](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-leaf-commit.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
