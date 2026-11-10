# pin-format-five

decision · `01M739GHYM7Z506MEACDWYCMFE`

## statement

The pin entry is a new log entry type in format 5: only pin entries carry v 5, so formats 1 to 4 stay byte-identical. A pin applies only while its change set is in flight, and readers built before format 5 refuse a log holding one with ErrNewerFormat.

## wording

the implementing agent's choice, accepted by the orchestrator acting as arbiter

## Links

### In

- decided-by [change-workflow-pin](../capability/rhiza/workflows/workflow-dispatch/change-workflow-pin.md)

## History

- `01M71G46Q5378CYYVEHFF2EMP3` Docs sign-off workflow: a repo workflow named docs (ship-change with signoff = owner, builders and fixers told to write with the humanize skill) builds doc change sets, and rhiza change new --workflow and a command for an existing draft pin a change set's workflow
