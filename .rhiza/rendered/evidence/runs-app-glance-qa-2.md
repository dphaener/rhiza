# runs-app-glance-qa-2

evidence · `01M6NC2354PNAV17NPZ9J8SYDS` · actor [claude-code-d7b0bf18](../actor/claude-code-d7b0bf18.md)

## result

PASS

## statement

Built HEAD (519eba52) to a scratch binary and ran it in the live repo: `rhiza tui app runs --once` and `rhiza run ls --json` gave identical bytes (cmp); in a pty, `rhiza tui app runs` listed exactly the 6 live runs of 20 (none finished) as title · current step · elapsed, elapsed ticked on refresh, q closed it with exit 0; `rhiza tui app health` → enter on the daemon row opened the same Runs panel and q closed it; `rhiza tui app live --once` exits 4 unknown app "live"; live-app-runs and live-app-focus are proposed retired in the slice. go test ./internal/tui/runsapp and the tui-runs-app and run-log testscripts pass.

## Links

### Out

- verifies [runs-app-glance](../capability/rhiza/runs-app-glance.md)

## History

- `01M6MYNJGY6Q735KB2ERT3V4R5` Runs app: a read-only glance at live runs and their current steps, rebuilt on the modal shell over run.ls, replacing internal/tui/liveapp and retiring the capabilities of the old live app it replaces
