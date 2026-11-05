# runs-app-glance-qa-1

evidence · `01M6NBA10AKQ71B21M9B6HFZGT` · actor [claude-code-840da7aa](../actor/claude-code-840da7aa.md)

## result

PASS

## statement

At 3a81c930, built rhiza with go build and ran it by hand against this repo's real runs (6 live, several finished). rhiza tui app runs --once was byte-identical (cmp) to rhiza run ls --json, and rhiza tui app live now exits 4 with unknown app "live". Interactive in tmux (120x30): an InfoPanel titled Runs listed exactly the 6 live runs, one line each with change set title, current step (review/review, build/build, wait-pr) and elapsed time; finished runs were not listed. Five seconds later the elapsed times had advanced with no key pressed, so it refreshes on its own. Arrow keys scrolled; q closed with exit 0. From tui app health, enter on the daemon row opened the Runs panel. Minor: there the health report's top three lines stayed drawn above the panel, which the old live app did not do.

## Links

### Out

- verifies [runs-app-glance](../capability/rhiza/runs-app-glance.md)

## History

- `01M6MYNJGY6Q735KB2ERT3V4R5` Runs app: a read-only glance at live runs and their current steps, rebuilt on the modal shell over run.ls, replacing internal/tui/liveapp and retiring the capabilities of the old live app it replaces
