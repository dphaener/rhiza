# build-ruling-9

decision · `01M5ZARTBVNX6Q40T6Y0WJF517`

## options

```json
[
  "Closed when the run ends",
  "Kept open after the run"
]
```

## picked

Closed when the run ends

## question

Are a failed or blocked session's panes closed when the run ends?

## session

c1b3cd7b-ff97-477d-b892-01a17e9f9184

## statement

Are a failed or blocked session's panes closed when the run ends? — Closed when the run ends. build-panes' "stays until the run ends" means the engine closes every remaining pane, failed or blocked sessions' included, when the run finishes.

## wording

the orchestrator's reading, confirmed in session c1b3cd7b-ff97-477d-b892-01a17e9f9184

## Links

### In

- decided-by [build-panes](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-panes.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
