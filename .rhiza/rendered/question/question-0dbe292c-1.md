# question-0dbe292c-1

question · `01M77VY0DCQKZF713XEZRQ9DRG`

## session

0dbe292c-3574-4c27-9bf7-01a1cfbd59e9

## statement

Change set 01M77VT6 (testscripts run inside a ship run write into ~/code/rhiza's main tree) serves no intent yet. The last run-repair bug fix (#484) served m4-runs-hold-up ('a run finishes without hand repair'), and the orchestrator had to move a stray file by hand here. Which intent should it serve? / How should the fix be shaped? Cause: workflows set GOTMPDIR={run_tmp} = <git common dir>/rhiza/runs/<id>/tmp, and go-internal testscript creates $WORK under $GOTMPDIR, so scripts that don't init their own repo reach ~/code/rhiza.

## Links

### In

- answers [decision-run-tmp-outside-repo](../decision/decision-run-tmp-outside-repo.md)

## History

- `01M77VT6TMYCDVMX61ZVA8CS4C` Fix Testscripts run inside a ship run write into the real repo's main tree: the run's GOTMPDIR lies under ~/code/rhiza/.git (.git/rhiza/checkouts/...), so a script's $WORK sits inside the real repo, and git/rhiza commands in scripts that don't init their own repo resolve to ~/code/rhiza. Evidence: run 01M77DWN (cs 01M763Z4, PR #480), builder session 91d8e5b2 ran 'go test ./cmd/rhiza -run TestScripts/' at 2026-11-12T09:14:48Z; at 09:15:17 an untracked .github/workflows/rhiza.yml ('Written by rhiza ci install', from cmd/rhiza/testdata/script/ci-install.txtar, internal/api/ci.go:23 workflowPath) appeared in ~/code/rhiza's main tree. Moved aside to ~/orchestrator/stray-logs/2026-11-12/rhiza.yml by the orchestrator. Untracked files on main block close-out pulls and break the 'nothing uncommitted on main' rule. The agent later worked around it with GOTMPDIR=scratchpad. Fix ideas: runs set GOTMPDIR/TMPDIR outside any git repo (e.g. under the session scratchpad or ~/.cache/rhiza/tmp); TestMain refuses to start when $WORK resolves inside a git work tree other than its own.
