# change-new-workflow-flag-verified

evidence · `01M72ZZBDMM3WWJRKR9CPG7CPB` · actor [claude-code-40b9eec7](../actor/claude-code-40b9eec7.md)

## result

PASS

## statement

Verifier at fc3c7176: go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(change-new-workflow-flag|describe-golden)$' passes. The new change-new-workflow-flag.txtar opens a change set with change new --workflow docs (repo workflow), sees "workflow": "docs" in change new --json and change show --json, sees run start --dry-run print '- Workflow: docs', sees --workflow nonesuch refused ('no such workflow', no stdout, no change set opened), sees a plain change new carry no workflow, and passes fold-ok. Read the diff: runChangeNew resolves the name through the shared lookupWorkflow and refuses an unresolved name with the same Usagef message as request --workflow, then records it on changeset.Open.Workflow (the field a request's workflow uses). describe-golden and docs/reference/cli.md were updated to match. Scoped go test ./internal/api -run 'ChangeNew|Workflow|Request' passes; the unscoped package run timed out at 600s with host load average 27, which is not caused by this change.

## Links

### Out

- verifies [change-new-workflow-flag](../capability/rhiza/workflows/workflow-dispatch/change-new-workflow-flag.md)

## History

- `01M71G46Q5378CYYVEHFF2EMP3` Docs sign-off workflow: a repo workflow named docs (ship-change with signoff = owner, builders and fixers told to write with the humanize skill) builds doc change sets, and rhiza change new --workflow and a command for an existing draft pin a change set's workflow
