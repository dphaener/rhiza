# build-ruling-6

decision · `01M5ZARSX9TBVG4JD0CKHA82J6`

## options

```json
[
  "Tests with each leaf",
  "TDD strictly"
]
```

## picked

Tests with each leaf

## question

What is the test strategy?

## session

c1b3cd7b-ff97-477d-b892-01a17e9f9184

## statement

What is the test strategy? — Tests with each leaf. Each task writes the tests its leaf's acceptance names (engine tests with a fake runner and clock, adapter tests, and testscripts with fake rex and claude), all hermetic, with package tests per task and make check before each commit.

## wording

picked by the orchestrator acting as arbiter

## Links

### In

- decided-by [build-workflow](../capability/rhiza/workflows/built-in-workflows/build-workflow.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
