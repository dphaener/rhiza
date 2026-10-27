# build-ruling-4

decision · `01M5ZARSKJE9V1QJ1K58HRHG0R`

## options

```json
[
  "Reason-only list",
  "Reuse plan finding kinds"
]
```

## picked

Reason-only list

## question

What shape do a verifier's findings take in `rhiza run findings --leaf <leaf>`?

## session

c1b3cd7b-ff97-477d-b892-01a17e9f9184

## statement

What shape do a verifier's findings take in `rhiza run findings --leaf <leaf>`? — Reason-only list. `rhiza run findings --leaf <leaf>` takes {"findings":[{"reason":"..."}]}, journaled as a findings entry carrying the leaf and attempt, and plan rounds ignore it; the leaf finding shape stays open to the check, severity, leaf and place fields change set 01M5WPTZ1V0JVVW8FPX0QYMRWJ (the review workflow) adds.

## wording

picked by the orchestrator acting as arbiter

## Links

### In

- decided-by [build-verifier-findings](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-verifier-findings.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
