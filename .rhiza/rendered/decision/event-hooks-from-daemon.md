# event-hooks-from-daemon

decision · `01M6BV2EV8GJG7S4MV9HF8HHSV`

## statement

The daemon, each poll, runs the event command (--event-command, or else git config rhiza.eventCommand) through sh -c in the repo root once per lifecycle move it has not handled, in log order, with the move as JSON on stdin and RHIZA_EVENT, RHIZA_EVENT_CHANGE, RHIZA_EVENT_ENTRY and RHIZA_EVENT_TITLE set; it records handled entry IDs in its state, seeding every existing move on first run so history is not replayed; a failing or timed-out (10 s) command is logged and not retried

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m17-orchestrator-picks](../source/m17-orchestrator-picks.md)

### In

- decided-by [event-hook-daemon](../capability/rhiza/hooks/model-event-hooks/event-hook-daemon.md)

## History

- `01M6BV1J7WC184NBNWQCX7C9WJ` Plan milestone 17: model events run your scripts, picked by the orchestrator under the owner's standing instruction
