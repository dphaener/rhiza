# live-app-runs-qa-1

evidence · `01M6993DFMZ1HHJFGH3BXDQEMF` · actor [claude-code-84d1030f](../actor/claude-code-84d1030f.md)

## result

FAIL

## statement

QA at 0f0c797: rhiza tui app live --once printed bytes identical to rhiza run ls --json (diff empty). Running rhiza tui app live in a 24-row Rex pane over this repo's 20 runs showed each run's title, state, step, leaves passed/conducted, checks and elapsed, and elapsed advanced by 3s between captures 4s apart with no key press, so the 3s reload works. But the view has no height handling: the daemon line, the two live runs at the top and the cursor mark were scrolled off, and only the oldest finished runs and the key help were visible. Finished runs' elapsed also keeps counting to now (27h+).

## Links

### Out

- verifies live-app-runs (retired)

## History

- `01M686RW277XSQ20MYCE580GG8` A live app shows the daemon and each live or recent run's progress, backed by a run.ls operation, and the roadmap detail refreshes its latest run
