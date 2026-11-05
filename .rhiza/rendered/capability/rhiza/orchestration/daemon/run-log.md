# run-log

capability · `01M6FESK5TDQ0J8S077V8BCZ12` · parent [daemon](../daemon.md)

## acceptance

A testscript: rhiza run log <run> prints the run's log file a page at a time with the paging flags ls takes, and exits 2 for a run with no log

## statement

rhiza run log <run> pages a headless run's driver log

## Links

### In

- verifies [run-log-verified](../../../../evidence/run-log-verified.md)
- verifies [run-log-qa-1](../../../../evidence/run-log-qa-1.md)
- verifies [run-log-verified-2](../../../../evidence/run-log-verified-2.md)
- verifies [run-log-qa-2](../../../../evidence/run-log-qa-2.md)
- verifies [pr-210-run-drivers-run-headless-as](../../../../evidence/pr-210-run-drivers-run-headless-as.md)
- verifies [run-log-runs-app-qa-2](../../../../evidence/run-log-runs-app-qa-2.md)
- verifies [pr-296-runs-app-on-the-modal-shell](../../../../evidence/pr-296-runs-app-on-the-modal-shell.md)

## History

- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
- `01M6MYNJGY6Q735KB2ERT3V4R5` Runs app: a read-only glance at live runs and their current steps, rebuilt on the modal shell over run.ls, replacing internal/tui/liveapp and retiring the capabilities of the old live app it replaces
