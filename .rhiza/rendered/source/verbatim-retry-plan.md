# verbatim-retry-plan

source · `01M7KWHE93J0DWFBZ47DPJ9Y1Z`

## summary

Autonomous plan of change set 01M76XT9 by the conductor of run 01M7KWCM5HATP75QHM7FE9J8R0; no owner was present. Research found: leafPrompt (internal/api/run_leaves.go) tells a retry builder only to read findings from rhiza run status --json; LeafState.AttemptFindings (internal/run/journal.go) already holds a verifier's journaled findings per attempt; rhiza run spawn writes a child's brief.md as the contract section plus the caller's brief (run.RenderSpawnBrief); the run journal has no notes; the build-retry topic already says inspect the working tree and address every finding as written. The plan: the leaves step's retry prompt quotes the earlier attempt's findings verbatim; a rhiza run note command journals one-line learnings; leaf prompts and spawned builder and verifier briefs carry the run's notes; the builder, verifier, build-retry and conduct topics say to write notes, quote findings verbatim and update landed nodes a leaf changes or raise a question.

## title

Verbatim retry and run notes, planned by the conductor

## Links

### In

- derived-from [run-notes-journaled](../decision/run-notes-journaled.md)
- derived-from [retry-prompt-quotes-findings](../capability/rhiza/workflows/built-in-workflows/build-workflow/retry-prompt-quotes-findings.md)

## History

- `01M76XT9E2N9ZYMMB8DMTMKA3Q` A leaf's retry prompt carries the verifier's findings verbatim and says to inspect the working tree and address every finding as written; builders and verifiers append one-line learnings to the run's notes, which later leaves' prompts carry; builders update landed nodes their leaf changes or raise a question. Evidence: 113 of 368 leaves needed a second attempt, 36 model-check majors were drift (prompt-gap recs 9-11).
