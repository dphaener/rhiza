# decision-run-tmp-outside-repo

decision · `01M77VZKS75A2ZFDRNFNYETS5E`

## reasoning

The question asks which intent and what shape. Intent: m4-runs-hold-up says a run finishes without hand repair; here the orchestrator had to move a stray .github/workflows/rhiza.yml off main by hand, and the previous run-repair bug fix (#484) served the same intent, so it is the clear fit. Shape: the cause is run.TmpDir (internal/run/store.go:45) placing run tmp at <git common dir>/rhiza/runs/<id>/tmp and ship-change/bugfix/plan-bug/docs workflows setting GOTMPDIR={run_tmp}; go-internal testscript makes $WORK under GOTMPDIR, so scripts that run git/rhiza without their own repo resolve to ~/code/rhiza. Moving run tmp out of the repo fixes every workflow and any other tool that uses {run_tmp} at once without editing the .rhiza/workflows files, unlike per-workflow overrides. The Setup guard is the change set's second fix idea and is cheap defense in depth: it catches any other way $WORK lands in a repo (a hand-set GOTMPDIR, a future workflow), turning a silent write to main into a clear test failure. Both ideas come from the change set's own intent text; I kept the run tmp lifecycle (made at step start, removed at run end) unchanged so nothing else in the engine moves.

## statement

Change set 01M77VT6 serves m4-runs-hold-up. Fix it at the source and guard the tests: (1) a run's {run_tmp} moves outside every git repository, to a per-user cache dir (os.UserCacheDir()/rhiza/run-tmp/<run id>, falling back to os.TempDir()/rhiza-run-tmp-<uid>/<run id>), still made as a step starts and removed when the run ends, so GOTMPDIR and anything else a workflow points at {run_tmp} no longer sit under <git common dir>; run.TmpDir/Store.TmpDir and the {run_tmp} docs change accordingly and existing workflows need no edit. (2) cmd/rhiza's testscript Setup fails a script whose $WORK resolves inside a git work tree or git dir (git rev-parse from $WORK before the script runs), naming the path and saying to set GOTMPDIR outside the repo, so no script can reach a real repo again. A testscript or unit test covers each: {run_tmp} resolves outside the repo's common dir, and Setup refuses a $WORK under a git repo.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-0dbe292c-1](../question/question-0dbe292c-1.md)
- derived-from [testscripts-in-run-plan](../source/testscripts-in-run-plan.md)

### In

- decided-by [run-tmp-outside-repo](../capability/rhiza/workflows/workflow-definitions/workflow-file/run-tmp-outside-repo.md)
- decided-by [testscripts-refuse-work-in-repo](../capability/rhiza/integrity/self-check/testscripts-refuse-work-in-repo.md)

## History

- `01M77VT6TMYCDVMX61ZVA8CS4C` Fix Testscripts run inside a ship run write into the real repo's main tree: the run's GOTMPDIR lies under ~/code/rhiza/.git (.git/rhiza/checkouts/...), so a script's $WORK sits inside the real repo, and git/rhiza commands in scripts that don't init their own repo resolve to ~/code/rhiza. Evidence: run 01M77DWN (cs 01M763Z4, PR #480), builder session 91d8e5b2 ran 'go test ./cmd/rhiza -run TestScripts/' at 2026-11-12T09:14:48Z; at 09:15:17 an untracked .github/workflows/rhiza.yml ('Written by rhiza ci install', from cmd/rhiza/testdata/script/ci-install.txtar, internal/api/ci.go:23 workflowPath) appeared in ~/code/rhiza's main tree. Moved aside to ~/orchestrator/stray-logs/2026-11-12/rhiza.yml by the orchestrator. Untracked files on main block close-out pulls and break the 'nothing uncommitted on main' rule. The agent later worked around it with GOTMPDIR=scratchpad. Fix ideas: runs set GOTMPDIR/TMPDIR outside any git repo (e.g. under the session scratchpad or ~/.cache/rhiza/tmp); TestMain refuses to start when $WORK resolves inside a git work tree other than its own.
