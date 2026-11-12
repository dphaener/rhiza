# testscripts-in-run-plan

source · `01M77W0KWA2PE04QA93PQ5A08H`

## quotes

```json
[]
```

## raw_sha256

31eba27ea3ec6886167412cbc96a7a789816d9e06d3f3c1840ff96a79b193e9e

## session

0dbe292c-3574-4c27-9bf7-01a1cfbd59e9

## summary

Autonomous plan run with no owner; the arbiter answered the planner's questions. Cause found: Rhiza's workflows set GOTMPDIR={run_tmp}, which is <git common dir>/rhiza/runs/<run id>/tmp, and go-internal testscript makes $WORK under $GOTMPDIR, so scripts that don't init their own repo reach ~/code/rhiza (ci-install wrote .github/workflows/rhiza.yml into main). The arbiter ruled: serve m4-runs-hold-up; move {run_tmp} out of every git repository to a per-user cache dir (os.UserCacheDir()/rhiza/run-tmp/<run id>, falling back to os.TempDir()/rhiza-run-tmp-<uid>/<run id>), still made at step start and removed at run end, with no workflow edits; and make cmd/rhiza's testscript Setup fail a script whose $WORK lies inside a git work tree or git dir, naming the path and saying to set GOTMPDIR outside the repo. Nothing left open.

## title

Plan: testscripts inside a run write into the real repo

## Links

### In

- derived-from [decision-run-tmp-outside-repo](../decision/decision-run-tmp-outside-repo.md)

## History

- `01M77VT6TMYCDVMX61ZVA8CS4C` Fix Testscripts run inside a ship run write into the real repo's main tree: the run's GOTMPDIR lies under ~/code/rhiza/.git (.git/rhiza/checkouts/...), so a script's $WORK sits inside the real repo, and git/rhiza commands in scripts that don't init their own repo resolve to ~/code/rhiza. Evidence: run 01M77DWN (cs 01M763Z4, PR #480), builder session 91d8e5b2 ran 'go test ./cmd/rhiza -run TestScripts/' at 2026-11-12T09:14:48Z; at 09:15:17 an untracked .github/workflows/rhiza.yml ('Written by rhiza ci install', from cmd/rhiza/testdata/script/ci-install.txtar, internal/api/ci.go:23 workflowPath) appeared in ~/code/rhiza's main tree. Moved aside to ~/orchestrator/stray-logs/2026-11-12/rhiza.yml by the orchestrator. Untracked files on main block close-out pulls and break the 'nothing uncommitted on main' rule. The agent later worked around it with GOTMPDIR=scratchpad. Fix ideas: runs set GOTMPDIR/TMPDIR outside any git repo (e.g. under the session scratchpad or ~/.cache/rhiza/tmp); TestMain refuses to start when $WORK resolves inside a git work tree other than its own.
