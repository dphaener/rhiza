# event-hook-daemon-verified

evidence · `01M6CAAC5Y8N9108JFMRQRT3MB` · actor [claude-code-3a233949](../actor/claude-code-3a233949.md)

## result

PASS

## statement

At 06603574 ran go test ./cmd/rhiza -run 'TestScripts/(daemon-event-hook|describe-golden|daemon-notify-hook|daemon-notify)$' -count=1: all PASS. daemon-event-hook.txtar uses a fake hook.sh appending stdin plus RHIZA_EVENT:RHIZA_EVENT_CHANGE:ENTRY:TITLE to a file and checks: first --once on a log already holding a propose writes nothing; after propose and accept one --once writes two lines, proposed then accepted, and a second --once writes none; a propose committed in a second clone and pulled fires exactly once; git config rhiza.eventCommand is used without the flag; a command exiting 1 is logged once as 'event command for entry <ULID>: exit status 1 ... boom' and the next poll does not rerun it; 'sleep 30' with a 200ms timeout is logged as 'event command for entry ... timed out after 200ms'. Regenerating describe-golden with -update left the tree clean (event_command is in it). Read the diff: daemonEvents sorts lifecycle moves by entry ID, seeds on first run, records each move before running so failures are not retried, and shares runHookCommand with the notify hook. make check exits 0.

## Links

### Out

- verifies [event-hook-daemon](../capability/rhiza/hooks/model-event-hooks/event-hook-daemon.md)

## History

- `01M6BV1MCF8DZV2GG1T852RBJM` The daemon runs an owner-configured command once per change set lifecycle event, including events pulled from other clones
