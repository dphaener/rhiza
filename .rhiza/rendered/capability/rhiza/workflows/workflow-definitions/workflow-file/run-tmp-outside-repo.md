# run-tmp-outside-repo

capability · `01M77W0Z8V8WBNPS6AZAG67R86` · parent [workflow-file](../workflow-file.md)

## acceptance

Unit tests in internal/run and internal/api (go test -trimpath ./internal/run/... ./internal/api/... -run for the run-tmp tests, e.g. gotmp_test.go and run_abandon_gotmp_test.go) show run.TmpDir resolves outside the repo's git common dir, a step sees {run_tmp} at that path with the folder present, and after the run ends or is abandoned only that run's folder is removed (never the shared run-tmp parent nor another run's folder); the {run_tmp} wording in internal/workflow/spec.go and the workflow-env docs names the new location

## statement

A run's {run_tmp} folder lies outside every git repository, at os.UserCacheDir()/rhiza/run-tmp/<run id> (or os.TempDir()/rhiza-run-tmp-<uid>/<run id> when there is no user cache dir), made as a step starts and removed when the run finishes, fails or is abandoned, so GOTMPDIR and anything else a workflow points at {run_tmp} never sits under the repo's git dir; existing workflows need no edit

## Links

### Out

- decided-by [decision-run-tmp-outside-repo](../../../../../decision/decision-run-tmp-outside-repo.md)

### In

- verifies [run-tmp-outside-repo-verified](../../../../../evidence/run-tmp-outside-repo-verified.md)
- verifies [run-tmp-outside-repo-qa-1](../../../../../evidence/run-tmp-outside-repo-qa-1.md)
- verifies [run-tmp-outside-repo-verified-2](../../../../../evidence/run-tmp-outside-repo-verified-2.md)
- verifies [run-tmp-outside-repo-qa-2](../../../../../evidence/run-tmp-outside-repo-qa-2.md)
- verifies [pr-488-fix-testscripts-run-inside-a](../../../../../evidence/pr-488-fix-testscripts-run-inside-a.md)

## History

- `01M77VT6TMYCDVMX61ZVA8CS4C` Fix Testscripts run inside a ship run write into the real repo's main tree: the run's GOTMPDIR lies under ~/code/rhiza/.git (.git/rhiza/checkouts/...), so a script's $WORK sits inside the real repo, and git/rhiza commands in scripts that don't init their own repo resolve to ~/code/rhiza. Evidence: run 01M77DWN (cs 01M763Z4, PR #480), builder session 91d8e5b2 ran 'go test ./cmd/rhiza -run TestScripts/' at 2026-11-12T09:14:48Z; at 09:15:17 an untracked .github/workflows/rhiza.yml ('Written by rhiza ci install', from cmd/rhiza/testdata/script/ci-install.txtar, internal/api/ci.go:23 workflowPath) appeared in ~/code/rhiza's main tree. Moved aside to ~/orchestrator/stray-logs/2026-11-12/rhiza.yml by the orchestrator. Untracked files on main block close-out pulls and break the 'nothing uncommitted on main' rule. The agent later worked around it with GOTMPDIR=scratchpad. Fix ideas: runs set GOTMPDIR/TMPDIR outside any git repo (e.g. under the session scratchpad or ~/.cache/rhiza/tmp); TestMain refuses to start when $WORK resolves inside a git work tree other than its own.
