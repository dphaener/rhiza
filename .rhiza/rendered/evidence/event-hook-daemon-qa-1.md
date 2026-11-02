# event-hook-daemon-qa-1

evidence · `01M6CAEB2FSTKTSB59DTYTN5JV` · actor [claude-code-68f2d78c](../actor/claude-code-68f2d78c.md)

## result

PASS

## statement

QA at fb5cc05c with a scratch build (go build -o $TMPDIR/rhiza-qa ./cmd/rhiza) in a fresh git repo, RHIZA_* unset, hook.sh appending stdin plus RHIZA_EVENT/CHANGE/ENTRY/TITLE and pwd to a file. rhiza daemon --permission-mode manual --once --event-command hook.sh on a log already holding a propose wrote nothing and recorded both existing entry IDs in .git/rhiza/daemon.state.json events. After accept, one poll wrote one line: JSON {event:accepted, change, entry, title (with embedded quotes escaped), at, actor} and matching env vars, run in the repo root; a second poll wrote none. Propose+accept of a second change set gave two lines in one poll, proposed then accepted. A propose committed in a second clone and pulled fired exactly once, and not again. With no flag, git config rhiza.eventCommand ran the hook on a withdraw. --event-command 'echo boom; exit 3' logged 'rhiza: daemon: event command for entry <ULID>: exit status 3: boom', the daemon exited 0, and the next poll did not rerun it. --event-command 'sleep 30' with no test override was killed after 10 s and logged 'timed out after 10s'.

## Links

### Out

- verifies [event-hook-daemon](../capability/rhiza/hooks/model-event-hooks/event-hook-daemon.md)

## History

- `01M6BV1MCF8DZV2GG1T852RBJM` The daemon runs an owner-configured command once per change set lifecycle event, including events pulled from other clones
