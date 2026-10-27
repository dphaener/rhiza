# build-signoff

decision · `01M5ZNVYW7BVEA0RXM60X80EEG`

## options

```json
[
  "Ship it + draft the follow-ups",
  "Ship it",
  "Changes needed"
]
```

## picked

Ship it + draft the follow-ups

## question

All four reviews approved and CI passes on PR #44. Do you approve shipping it?

## session

c1b3cd7b-ff97-477d-b892-01a17e9f9184

## statement

All four reviews approved and CI passes on PR #44. Do you approve shipping it? — Ship it + draft the follow-ups. The build workflow's change set is signed off for shipping, and four noticed issues outside its intent are recorded as discovered-from draft change sets, never built in this run: builder actor records left untracked (01M5ZNV3K796MKWEZM3TBQXGZ8), concurrent run findings sharing an ID (01M5ZNV40RB8VYBB0A3CPXCNXZ), a second live run on one change set (01M5ZNV4E4AJR52EX631180DMS), and the pre-existing TestGuideInstallConcurrent flake, to be fixed at its cause without skipping or loosening the test (01M5ZNVHN2YNK6DDVYY8GK5CEQ).

## wording

picked by the orchestrator acting as arbiter

## Links

### In

- decided-by [build-workflow](../capability/rhiza/workflows/built-in-workflows/build-workflow.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
