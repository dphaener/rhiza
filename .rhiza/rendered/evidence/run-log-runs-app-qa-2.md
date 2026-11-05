# run-log-runs-app-qa-2

evidence · `01M6NC29FENXQ80B7REWHAV0S9` · actor [claude-code-d7b0bf18](../actor/claude-code-d7b0bf18.md)

## result

PASS

## statement

Built HEAD (519eba52) and ran `rhiza run log 01M6NA07P221KN3WS9E47T94G3 --limit 3` in the live repo: it printed the first 3 driver-log lines and a `next: ... --cursor 3` line, exit 0; the run-log testscript (incl. exit 2 for a run with no log) passes.

## Links

### Out

- verifies [run-log](../capability/rhiza/orchestration/daemon/run-log.md)

## History

- `01M6MYNJGY6Q735KB2ERT3V4R5` Runs app: a read-only glance at live runs and their current steps, rebuilt on the modal shell over run.ls, replacing internal/tui/liveapp and retiring the capabilities of the old live app it replaces
