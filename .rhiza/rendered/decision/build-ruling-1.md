# build-ruling-1

decision · `01M5ZARS4Y57P2RCNQ25X4ZBKX`

## options

```json
[
  "One polling loop",
  "Goroutine per session"
]
```

## picked

One polling loop

## question

How should the leaves step conduct a builder and parallel verifiers inside the deterministic engine?

## session

c1b3cd7b-ff97-477d-b892-01a17e9f9184

## statement

How should the leaves step conduct a builder and parallel verifiers inside the deterministic engine? — One polling loop. The leaves step is one polling loop on the engine's goroutine over every live leaf session, one builder and any number of verifiers, journaling each change; no goroutine per session.

## wording

picked by the orchestrator acting as arbiter

## Links

### In

- decided-by [build-leaves-step](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-leaves-step.md)
- decided-by [build-leaf-order](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-leaf-order.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
