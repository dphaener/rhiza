# live-app-runs-qa-2

evidence · `01M69ABJHKR5PH6R1W7P7QB05X` · actor [claude-code-65ab0a4c](../actor/claude-code-65ab0a4c.md)

## result

FAIL

## statement

QA cycle 2 at 1d0f3f4, built binary in a real 80x24 Rex pane over this repo's 20 runs: rhiza tui app live showed the daemon line, 7 runs and the key help; 9 j presses scrolled so the cursor's run (10th) sat at the bottom with the daemon line kept; 15 more j stopped on the 20th run; resizing the pane to 10 rows re-fit to 2 runs with the cursor kept; 25 k returned to the top. The live run 01M696N43 has a review/review fixer entry for live-app-runs in run.status and the app read 'leaves 4/4' (the pre-fix build would read 4/5). But each reload now calls change.show per change set (~0.66s each): tui app live --once took 13.2s against 0.7s for the pre-fix build 2d2519d^, and since Reload runs inside Update every 3s, keys lagged: q took 7s to quit against 0s for the pre-fix build.

## Links

### Out

- verifies live-app-runs (retired)

## History

- `01M686RW277XSQ20MYCE580GG8` A live app shows the daemon and each live or recent run's progress, backed by a run.ls operation, and the roadmap detail refreshes its latest run
