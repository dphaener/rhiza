# run-tmp-outside-repo-qa-1

evidence · `01M77XMRTNCX6FHAV031SF6M9T` · actor [claude-code-3fd2acf3](../actor/claude-code-3fd2acf3.md)

## result

PASS

## statement

QA at 1d8815c86 with a binary built from the branch: in a scratch Rhiza root and git repo, a workflow with [env] GOTMPDIR = "{run_tmp}" and one command step, run via rhiza run start --foreground. The run succeeded; the step saw GOTMPDIR=~/.cache/rhiza/run-tmp/<run id>, the folder existed during the step, git rev-parse from it said 'not a git repository', and after the run the folder was gone while a sibling run-tmp folder (R1) and the shared parent stayed. Nothing was made under the repo's .git/rhiza/runs beyond the run journal.

## Links

### Out

- verifies [run-tmp-outside-repo](../capability/rhiza/workflows/workflow-definitions/workflow-file/run-tmp-outside-repo.md)

## History

- `01M77VT6TMYCDVMX61ZVA8CS4C` Fix Testscripts run inside a ship run write into the real repo's main tree: the run's GOTMPDIR lies under ~/code/rhiza/.git (.git/rhiza/checkouts/...), so a script's $WORK sits inside the real repo, and git/rhiza commands in scripts that don't init their own repo resolve to ~/code/rhiza. Evidence: run 01M77DWN (cs 01M763Z4, PR #480), builder session 91d8e5b2 ran 'go test ./cmd/rhiza -run TestScripts/' at 2026-11-12T09:14:48Z; at 09:15:17 an untracked .github/workflows/rhiza.yml ('Written by rhiza ci install', from cmd/rhiza/testdata/script/ci-install.txtar, internal/api/ci.go:23 workflowPath) appeared in ~/code/rhiza's main tree. Moved aside to ~/orchestrator/stray-logs/2026-11-12/rhiza.yml by the orchestrator. Untracked files on main block close-out pulls and break the 'nothing uncommitted on main' rule. The agent later worked around it with GOTMPDIR=scratchpad. Fix ideas: runs set GOTMPDIR/TMPDIR outside any git repo (e.g. under the session scratchpad or ~/.cache/rhiza/tmp); TestMain refuses to start when $WORK resolves inside a git work tree other than its own.
