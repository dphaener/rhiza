# build-ruling-2

decision · `01M5ZARS9RT9KNMCK8YMQ6KQF1`

## options

```json
[
  "RHIZA_ROOT env",
  "--root global flag"
]
```

## picked

RHIZA_ROOT env

## question

How do a verifier's rhiza writes reach the build worktree's log from its detached checkout?

## session

c1b3cd7b-ff97-477d-b892-01a17e9f9184

## statement

How do a verifier's rhiza writes reach the build worktree's log from its detached checkout? — RHIZA_ROOT env. A verifier's rhiza writes reach the build worktree's log through RHIZA_ROOT: the CLI takes its root from it when set, and the engine sets it only for verifier sessions.

## wording

picked by the orchestrator acting as arbiter

## Links

### In

- decided-by [build-verifier-checkout](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-verifier-checkout.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
