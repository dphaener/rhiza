# testscripts-refuse-work-in-repo-qa-2

evidence · `01M7805MXWNR7HPY94HFRXP106` · actor [claude-code-d38b579c](../actor/claude-code-d38b579c.md)

## result

PASS

## statement

QA cycle 2 at c3bc52840: with this run's own GOTMPDIR (~/code/rhiza/.git/rhiza/runs/<run>/tmp), go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/ci-install$' fails before the script runs ('inside a git dir: set GOTMPDIR outside the repository') and ~/code/rhiza/.github/workflows holds no rhiza.yml; with GOTMPDIR in the session scratchpad, TestCheckWorkOutsideRepo and TestScripts/ci-install pass.

## Links

### Out

- verifies [testscripts-refuse-work-in-repo](../capability/rhiza/integrity/self-check/testscripts-refuse-work-in-repo.md)

## History

- `01M77VT6TMYCDVMX61ZVA8CS4C` Fix Testscripts run inside a ship run write into the real repo's main tree: the run's GOTMPDIR lies under ~/code/rhiza/.git (.git/rhiza/checkouts/...), so a script's $WORK sits inside the real repo, and git/rhiza commands in scripts that don't init their own repo resolve to ~/code/rhiza. Evidence: run 01M77DWN (cs 01M763Z4, PR #480), builder session 91d8e5b2 ran 'go test ./cmd/rhiza -run TestScripts/' at 2026-11-12T09:14:48Z; at 09:15:17 an untracked .github/workflows/rhiza.yml ('Written by rhiza ci install', from cmd/rhiza/testdata/script/ci-install.txtar, internal/api/ci.go:23 workflowPath) appeared in ~/code/rhiza's main tree. Moved aside to ~/orchestrator/stray-logs/2026-11-12/rhiza.yml by the orchestrator. Untracked files on main block close-out pulls and break the 'nothing uncommitted on main' rule. The agent later worked around it with GOTMPDIR=scratchpad. Fix ideas: runs set GOTMPDIR/TMPDIR outside any git repo (e.g. under the session scratchpad or ~/.cache/rhiza/tmp); TestMain refuses to start when $WORK resolves inside a git work tree other than its own.
