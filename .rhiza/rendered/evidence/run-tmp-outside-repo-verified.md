# run-tmp-outside-repo-verified

evidence · `01M77WRJGPZJQBFVEHA0D4BC81` · actor [claude-code-1cff2f76](../actor/claude-code-1cff2f76.md)

## result

PASS

## statement

Verifier at de929a3bf: read the diff (run.TmpDir now os.UserCacheDir()/rhiza/run-tmp/<id>, falling back to os.TempDir()/rhiza-run-tmp-<uid>/<id>; engine and run.abandon remove only TmpDir(id), not its parent; session_run and watchdog callers use the same TmpDir; spec.go env wording and docs/reference/workflow-file.md name the new location). Ran go vet on internal/run, internal/api, internal/workflow (clean); go test -trimpath -count=1 ./internal/run/... -run 'RunTmp|TmpDir' -v: TestStepsGetEnvRunTmp (FOO={run_tmp} equals <XDG_CACHE_HOME>/rhiza/run-tmp/<id>, folder present when command, wait and agent steps run, gone after), TestRunTmpGoneAfterFailedRun, TestNoEnvNoRunTmp, TestTmpDirFallsBackToOSTemp, TestRunTmpRemovalSparesOthers all PASS; go test -trimpath -count=1 ./internal/api/... -run GoTmp: TestAbandonRemovesRunGoTmp PASS (own folder removed, other run's folder and shared parent kept); full ./internal/run/... and ./internal/workflow/... pass; gofmt clean.

## Links

### Out

- verifies [run-tmp-outside-repo](../capability/rhiza/workflows/workflow-definitions/workflow-file/run-tmp-outside-repo.md)

## History

- `01M77VT6TMYCDVMX61ZVA8CS4C` Fix Testscripts run inside a ship run write into the real repo's main tree: the run's GOTMPDIR lies under ~/code/rhiza/.git (.git/rhiza/checkouts/...), so a script's $WORK sits inside the real repo, and git/rhiza commands in scripts that don't init their own repo resolve to ~/code/rhiza. Evidence: run 01M77DWN (cs 01M763Z4, PR #480), builder session 91d8e5b2 ran 'go test ./cmd/rhiza -run TestScripts/' at 2026-11-12T09:14:48Z; at 09:15:17 an untracked .github/workflows/rhiza.yml ('Written by rhiza ci install', from cmd/rhiza/testdata/script/ci-install.txtar, internal/api/ci.go:23 workflowPath) appeared in ~/code/rhiza's main tree. Moved aside to ~/orchestrator/stray-logs/2026-11-12/rhiza.yml by the orchestrator. Untracked files on main block close-out pulls and break the 'nothing uncommitted on main' rule. The agent later worked around it with GOTMPDIR=scratchpad. Fix ideas: runs set GOTMPDIR/TMPDIR outside any git repo (e.g. under the session scratchpad or ~/.cache/rhiza/tmp); TestMain refuses to start when $WORK resolves inside a git work tree other than its own.
