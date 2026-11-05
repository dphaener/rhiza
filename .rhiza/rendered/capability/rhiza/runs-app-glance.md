# runs-app-glance

capability · `01M6MZBJWTAMCNY0T5ZBANQJNN` · parent [rhiza](../rhiza.md)

## acceptance

Go tests in the runs app's package with a fake runner: one row per live run with its title and current step, a finished run is not listed, a refresh tick updates the rows, no live runs shows one empty-state line, and the runner is only ever called with read operations; the runs app's testscript compares rhiza tui app runs --once with rhiza run ls --json byte for byte.

## statement

The Runs app is a read-only rexuikit List or InfoPanel of the live runs, each with its change set title, current step and elapsed time, refreshing on its own; it writes nothing, and --once prints what rhiza run ls --json prints.

## Links

### Out

- decided-by [tui-five-apps](../../decision/tui-five-apps.md)
- decided-by [tui-rexuikit-components](../../decision/tui-rexuikit-components.md)

### In

- verifies [runs-app-glance-verified](../../evidence/runs-app-glance-verified.md)
- verifies [runs-app-glance-qa-1](../../evidence/runs-app-glance-qa-1.md)
- verifies [runs-app-glance-verified-2](../../evidence/runs-app-glance-verified-2.md)
- verifies [runs-app-glance-qa-2](../../evidence/runs-app-glance-qa-2.md)
- verifies [pr-296-runs-app-on-the-modal-shell](../../evidence/pr-296-runs-app-on-the-modal-shell.md)

## History

- `01M6MYNJGY6Q735KB2ERT3V4R5` Runs app: a read-only glance at live runs and their current steps, rebuilt on the modal shell over run.ls, replacing internal/tui/liveapp and retiring the capabilities of the old live app it replaces
