# live-app-runs-verified-5

evidence · `01M69AQ7MHQE9211E4QBEWA55H` · actor [claude-code-508954ac](../actor/claude-code-508954ac.md)

## result

PASS

## statement

At 2254fea, go test -count=1 -v ./internal/tui/liveapp passed all 8 tests (also -race -count=3). Over a fake runner with canned daemon.status, run.ls, run.status and change.show and an injected clock, TestViewShowsRunProgress asserts the title, 'step review · leaves 2/4 · checks 1 passed, 1 failed, 1 running · elapsed 5m07s', and TestTickUpdatesViewKeepingCursorOnRun asserts that a tick with changed results updates the view (Daemon: not running, new Third run, leaves 3/4 then 3/3 after r) with the cursor kept on R1, and that change.show runs off the update loop. go test ./cmd/rhiza -run TestScripts/tui-live-app passed: rhiza tui app live --once cmp-matches rhiza run ls --json, both empty and with one run. TestEmptyAndRegistered finds the app in launcher.Apps(). I read the diff: leaf totals now come through a LeavesMsg command, and visible() reserves the blank line. make check: fmt, vet and tests pass; self-check fails only on 18 stale .rhiza/rendered files from run model writes outside this commit, which touches no .rhiza file.

## Links

### Out

- verifies live-app-runs (retired)

## History

- `01M686RW277XSQ20MYCE580GG8` A live app shows the daemon and each live or recent run's progress, backed by a run.ls operation, and the roadmap detail refreshes its latest run
