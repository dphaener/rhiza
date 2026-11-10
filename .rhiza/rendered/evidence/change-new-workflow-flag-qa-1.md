# change-new-workflow-flag-qa-1

evidence · `01M7398TVBV2XE0WEN8SY3YBTC` · actor [claude-code-32c9dfde](../actor/claude-code-32c9dfde.md)

## result

PASS

## statement

QA at 0e35c5a8 with the branch binary in a scratch Rhiza root and repo carrying this repo's docs workflows: change new 'Docs pass' --workflow docs --json printed "workflow": "docs" and change show --json showed it; run start --dry-run --permission-mode manual on it printed '- Workflow: docs' and '- Autonomy signoff: owner from workflow docs'; change new --workflow nonesuch exited 2 with 'no such workflow' and opened no change set; a plain change new carried no workflow; rhiza check passed after render (open entry stays format 4).

## Links

### Out

- verifies [change-new-workflow-flag](../capability/rhiza/workflows/workflow-dispatch/change-new-workflow-flag.md)

## History

- `01M71G46Q5378CYYVEHFF2EMP3` Docs sign-off workflow: a repo workflow named docs (ship-change with signoff = owner, builders and fixers told to write with the humanize skill) builds doc change sets, and rhiza change new --workflow and a command for an existing draft pin a change set's workflow
