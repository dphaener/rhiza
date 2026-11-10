# change-workflow-pin-qa-1

evidence · `01M7398TZDSG25VVSBC7CS31P0` · actor [claude-code-32c9dfde](../actor/claude-code-32c9dfde.md)

## result

PASS

## statement

QA at 0e35c5a8 with the branch binary in a scratch repo: change workflow docs --dry-run said 'Would pin' and wrote nothing; change workflow docs pinned it (change show --json "workflow": "docs", run start --dry-run '- Workflow: docs'); change workflow --clear --json reported previous docs, change show lost the workflow and run start --dry-run fell back to ship-change (build/build); nonesuch, no name, and name plus --clear each exited 2 with clear messages; after propose/accept/land, change workflow --clear --change <it> exited 3 'is landed'; --json and api call change.workflow printed identical bytes; only pin entries are format 5; rhiza check passed.

## Links

### Out

- verifies [change-workflow-pin](../capability/rhiza/workflows/workflow-dispatch/change-workflow-pin.md)

## History

- `01M71G46Q5378CYYVEHFF2EMP3` Docs sign-off workflow: a repo workflow named docs (ship-change with signoff = owner, builders and fixers told to write with the humanize skill) builds doc change sets, and rhiza change new --workflow and a command for an existing draft pin a change set's workflow
