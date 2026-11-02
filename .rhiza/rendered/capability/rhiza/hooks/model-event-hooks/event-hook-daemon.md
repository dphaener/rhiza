# event-hook-daemon

capability · `01M6BV2VY4YA51SRD25ZF3S09M` · parent [model-event-hooks](../model-event-hooks.md)

## acceptance

daemon-event-hook.txtar with a fake command appending stdin and RHIZA_EVENT:RHIZA_EVENT_CHANGE to a file: the first daemon --once on an existing log writes nothing; after propose and accept one --once writes exactly two lines in log order and a second writes none; a move committed in a second clone and pulled fires once; git config rhiza.eventCommand is used when the flag is absent; a failing or timed-out command is logged with its entry ID, the poll succeeds and it is not retried; describe-golden.txtar is regenerated

## statement

The daemon runs the user's event command once per lifecycle move it has not handled, in log order, with the move as JSON on stdin and its event, change set, entry and title in RHIZA_EVENT_* variables, including moves pulled from other clones, never replaying history on its first run

## Links

### Out

- decided-by [event-hooks-from-daemon](../../../../decision/event-hooks-from-daemon.md)

### In

- verifies [event-hook-daemon-verified](../../../../evidence/event-hook-daemon-verified.md)
- verifies [event-hook-daemon-qa-1](../../../../evidence/event-hook-daemon-qa-1.md)
- verifies [pr-182-lifecycle-event-hooks](../../../../evidence/pr-182-lifecycle-event-hooks.md)

## History

- `01M6BV1MCF8DZV2GG1T852RBJM` The daemon runs an owner-configured command once per change set lifecycle event, including events pulled from other clones
