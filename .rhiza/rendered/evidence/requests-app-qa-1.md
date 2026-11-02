# requests-app-qa-1

evidence · `01M6DEWRKK3S13DWM0BKXEHX7M` · actor [claude-code-50960276](../actor/claude-code-50960276.md)

## result

PASS

## statement

QA at 1916b724 with a scratch build (go build ./cmd/rhiza) in a scratch git repo: `rhiza tui app requests --once` was byte-equal (cmp) to `rhiza request ls --json` empty and with a bug, a feature and a milestone filed. Drove the interactive app in a 120x40 pty: b + 'crash on save' + enter queued a bug and moved the cursor to it ('Queued bug request ...'); n + enter on empty text kept the entry open with 'A request needs its text'; tab twice cycled feature->bug->milestone and 'beta  release ' filed as milestone 'beta release'; esc cancelled typed text; enter on a queued request said it has no run yet; with a request hand-marked with a dead run the row showed 'planning · run <id>' and enter said 'Run <id> is not live'; q exited 0; request ls showed all five queued. Against real Rex, `rhiza run open --run <this review run>` focused its live tab (opened true, change filled in); --run with a change set, or with neither, exits 2. `rhiza rex install` with REX_CONFIG in a scratch dir and a fake rex wrote rhiza_app_requests running `tui app requests` on key q; diffed against the merge-base build's key table, the only change is the added ["requests"] = "q" and no key is shared by two apps. tui-requests-app and rex-install scripts and requestsapp/api tests pass.

## Links

### Out

- verifies [requests-app](../capability/rhiza/orchestration/daemon/plan-requests/requests-app.md)

## History

- `01M6CRZRVGZBRB51YFPXHKZA75` A requests app files features, bugs and milestones and shows each request's progress
