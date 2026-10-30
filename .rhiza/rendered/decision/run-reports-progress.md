# run-reports-progress

decision · `01M5WT1X52NSSH59BY5KP418D0`

## quote

I like the idea of having what step the run is on and which leaf, etc. Let's add it as a leaf.

## quote_where

milestone-3 conversation, on OSC 7501 program status

## statement

A run reports its progress to Rex through OSC 7501 program status: the run is the root record, with its current step and progress as leaves done out of all; each live leaf is a child record with a hierarchical ID, the builder or verifier working or blocked; and the run ends done or error.

## wording

the orchestrator's proposal in the milestone-3 conversation, accepted by the owner

## Links

### Out

- derived-from [milestone-3-conversation](../source/milestone-3-conversation.md)

### In

- decided-by [run-program-status](../capability/rhiza/workflows/workflow-runs/run-program-status.md)
- decided-by [review-fixers-in-rex](../capability/rhiza/workflows/workflow-runs/run-program-status/review-fixers-in-rex.md)
- decided-by [review-code-check-in-rex](../capability/rhiza/workflows/workflow-runs/run-program-status/review-code-check-in-rex.md)
- decided-by [review-questions-in-rex](../capability/rhiza/workflows/workflow-runs/run-program-status/review-questions-in-rex.md)
- supersedes [run-ends-blocked-on-review-question](run-ends-blocked-on-review-question.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
