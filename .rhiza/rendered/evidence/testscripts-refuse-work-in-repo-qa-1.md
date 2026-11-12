# testscripts-refuse-work-in-repo-qa-1

evidence · `01M77XMRX1767S30VTFZC0NZE9` · actor [claude-code-3fd2acf3](../actor/claude-code-3fd2acf3.md)

## result

PASS

## statement

QA at 1d8815c86: go test -trimpath ./cmd/rhiza -run TestScripts/ci-install$ with GOTMPDIR set to a scratch repo's .git/rhiza/tmp failed before the script ran ('inside a git dir: set GOTMPDIR outside the repository'); with GOTMPDIR set to a folder in that repo's work tree it failed naming the work tree; neither run wrote .github/workflows/rhiza.yml or anything else into the repo (git status clean). With GOTMPDIR outside any repo ci-install passed, and TestCheckWorkOutsideRepo passes.

## Links

### Out

- verifies [testscripts-refuse-work-in-repo](../capability/rhiza/integrity/self-check/testscripts-refuse-work-in-repo.md)

## History

- `01M77VT6TMYCDVMX61ZVA8CS4C` Fix Testscripts run inside a ship run write into the real repo's main tree: the run's GOTMPDIR lies under ~/code/rhiza/.git (.git/rhiza/checkouts/...), so a script's $WORK sits inside the real repo, and git/rhiza commands in scripts that don't init their own repo resolve to ~/code/rhiza. Evidence: run 01M77DWN (cs 01M763Z4, PR #480), builder session 91d8e5b2 ran 'go test ./cmd/rhiza -run TestScripts/' at 2026-11-12T09:14:48Z; at 09:15:17 an untracked .github/workflows/rhiza.yml ('Written by rhiza ci install', from cmd/rhiza/testdata/script/ci-install.txtar, internal/api/ci.go:23 workflowPath) appeared in ~/code/rhiza's main tree. Moved aside to ~/orchestrator/stray-logs/2026-11-12/rhiza.yml by the orchestrator. Untracked files on main block close-out pulls and break the 'nothing uncommitted on main' rule. The agent later worked around it with GOTMPDIR=scratchpad. Fix ideas: runs set GOTMPDIR/TMPDIR outside any git repo (e.g. under the session scratchpad or ~/.cache/rhiza/tmp); TestMain refuses to start when $WORK resolves inside a git work tree other than its own.
