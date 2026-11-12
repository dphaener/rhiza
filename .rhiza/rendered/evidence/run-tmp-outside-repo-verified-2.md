# run-tmp-outside-repo-verified-2

evidence · `01M77ZY3PB94NDB8S0DNJW601G` · actor [claude-code-379896c8](../actor/claude-code-379896c8.md)

## result

PASS

## statement

Verifier, attempt 2, at 51a44876d (GOTMPDIR outside the repo). The two confirmed review findings are fixed. (1) go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/workflow-env$' passes; the script now expects FOO=<...>/.cache/rhiza/run-tmp/<run id>. (2) newSessionFixture sets XDG_CACHE_HOME to t.TempDir(). go test -trimpath -count=1 ./internal/api -run 'TestSessionRunWorkflowEnv$|Abandon|GoTmp' -v passes on two back-to-back runs (TestSessionRunWorkflowEnv, TestAbandonRemovesRunGoTmp, TestUniqueSlugShipTargetAndAbandon), and the host's ~/.cache/rhiza/run-tmp stays empty afterwards. go test -trimpath -count=1 ./internal/run/... passes. A verbose run of the whole internal/api package passed all 628 tests. The package was still marked FAIL because the root-state guard saw ~/code/.rhiza-root/usage.live.json change during the 14-minute run on the shared host. That is the live root daemon's file and is not touched by this diff, which only changes a test fixture env and one txtar.

## Links

### Out

- verifies [run-tmp-outside-repo](../capability/rhiza/workflows/workflow-definitions/workflow-file/run-tmp-outside-repo.md)

## History

- `01M77VT6TMYCDVMX61ZVA8CS4C` Fix Testscripts run inside a ship run write into the real repo's main tree: the run's GOTMPDIR lies under ~/code/rhiza/.git (.git/rhiza/checkouts/...), so a script's $WORK sits inside the real repo, and git/rhiza commands in scripts that don't init their own repo resolve to ~/code/rhiza. Evidence: run 01M77DWN (cs 01M763Z4, PR #480), builder session 91d8e5b2 ran 'go test ./cmd/rhiza -run TestScripts/' at 2026-11-12T09:14:48Z; at 09:15:17 an untracked .github/workflows/rhiza.yml ('Written by rhiza ci install', from cmd/rhiza/testdata/script/ci-install.txtar, internal/api/ci.go:23 workflowPath) appeared in ~/code/rhiza's main tree. Moved aside to ~/orchestrator/stray-logs/2026-11-12/rhiza.yml by the orchestrator. Untracked files on main block close-out pulls and break the 'nothing uncommitted on main' rule. The agent later worked around it with GOTMPDIR=scratchpad. Fix ideas: runs set GOTMPDIR/TMPDIR outside any git repo (e.g. under the session scratchpad or ~/.cache/rhiza/tmp); TestMain refuses to start when $WORK resolves inside a git work tree other than its own.
