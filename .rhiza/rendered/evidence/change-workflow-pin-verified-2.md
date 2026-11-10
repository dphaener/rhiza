# change-workflow-pin-verified-2

evidence · `01M739PRTK7PPB0JWBP7TBMC09` · actor [claude-code-b7639b9d](../actor/claude-code-b7639b9d.md)

## result

PASS

## statement

At bb5b5ce8 (leaf commits f82e440f and bb5b5ce8): go test -trimpath ./cmd/rhiza -run 'TestScripts/(change-workflow-pin|describe-golden)$' -count=1 -v passed both scripts. change-workflow-pin.txtar pins a draft change set to docs, then change show --json shows "workflow": "docs" and run start --dry-run prints '- Workflow: docs'. It clears with --clear (previous docs, no workflow, dry run falls back to the dispatch rules), refuses nonesuch with 'no such workflow' and a landed change set with 'landed', and passes fold-ok after each step. go test -trimpath ./internal/api -run TestChangeWorkflow passed, as did internal/changeset, entry, store and render. I read the diff: change.workflow validates the name via lookupWorkflow, refuses terminal states and live runs (refuseLiveRun), and appends a format-5 pin entry. bb5b5ce8 adds Decision pin-format-five, decided-by link. make self-check has only stale rendered actor/index files from this run's own session actors, which no leaf commit causes.

## Links

### Out

- verifies [change-workflow-pin](../capability/rhiza/workflows/workflow-dispatch/change-workflow-pin.md)

## History

- `01M71G46Q5378CYYVEHFF2EMP3` Docs sign-off workflow: a repo workflow named docs (ship-change with signoff = owner, builders and fixers told to write with the humanize skill) builds doc change sets, and rhiza change new --workflow and a command for an existing draft pin a change set's workflow
