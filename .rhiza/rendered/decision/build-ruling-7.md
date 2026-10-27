# build-ruling-7

decision · `01M5ZART24K8H2ANBT4AKCKFVJ`

## options

```json
[
  "On session done or exit",
  "On idle too"
]
```

## picked

On session done or exit

## question

When does a builder or verifier session count as finished?

## session

c1b3cd7b-ff97-477d-b892-01a17e9f9184

## statement

When does a builder or verifier session count as finished? — On session done or exit. A builder or verifier attempt is judged once its agent runs rhiza session done or its process exits, never on idle alone.

## wording

picked by the orchestrator acting as arbiter

## Links

### In

- decided-by [build-leaves-step](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-leaves-step.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
