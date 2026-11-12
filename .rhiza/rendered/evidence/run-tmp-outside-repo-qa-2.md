# run-tmp-outside-repo-qa-2

evidence · `01M7805MWNRKYY3DCJS629QT52` · actor [claude-code-d38b579c](../actor/claude-code-d38b579c.md)

## result

PASS

## statement

QA cycle 2 at c3bc52840 with a binary built from the branch, real ~/.cache. In a scratch Rhiza root and repo, a workflow with [env] GOTMPDIR = "{run_tmp}": the step saw GOTMPDIR=~/.cache/rhiza/run-tmp/<run id>, present during the step, git rev-parse there said 'not a git repository', go env GOTMPDIR agreed, nothing under the repo's .git/rhiza/runs/*/tmp; after the run its folder was gone and a sibling folder and the parent stayed. A second workflow step ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(ci-install|workflow-env)$' with GOTMPDIR={run_tmp}: ok, and git status of the checkout showed no new file (no .github/workflows/rhiza.yml). A third run was killed -9 mid-step, then rhiza run abandon removed only its own run-tmp folder, keeping the sibling and the run-tmp parent. Scoped tests with GOTMPDIR in the scratchpad pass: TestScripts/workflow-env, ./internal/run/..., ./internal/api -run 'TestSessionRunWorkflowEnv$|Abandon|GoTmp'; the host's ~/.cache/rhiza/run-tmp was empty before and after.

## Links

### Out

- verifies [run-tmp-outside-repo](../capability/rhiza/workflows/workflow-definitions/workflow-file/run-tmp-outside-repo.md)

## History

- `01M77VT6TMYCDVMX61ZVA8CS4C` Fix Testscripts run inside a ship run write into the real repo's main tree: the run's GOTMPDIR lies under ~/code/rhiza/.git (.git/rhiza/checkouts/...), so a script's $WORK sits inside the real repo, and git/rhiza commands in scripts that don't init their own repo resolve to ~/code/rhiza. Evidence: run 01M77DWN (cs 01M763Z4, PR #480), builder session 91d8e5b2 ran 'go test ./cmd/rhiza -run TestScripts/' at 2026-11-12T09:14:48Z; at 09:15:17 an untracked .github/workflows/rhiza.yml ('Written by rhiza ci install', from cmd/rhiza/testdata/script/ci-install.txtar, internal/api/ci.go:23 workflowPath) appeared in ~/code/rhiza's main tree. Moved aside to ~/orchestrator/stray-logs/2026-11-12/rhiza.yml by the orchestrator. Untracked files on main block close-out pulls and break the 'nothing uncommitted on main' rule. The agent later worked around it with GOTMPDIR=scratchpad. Fix ideas: runs set GOTMPDIR/TMPDIR outside any git repo (e.g. under the session scratchpad or ~/.cache/rhiza/tmp); TestMain refuses to start when $WORK resolves inside a git work tree other than its own.
