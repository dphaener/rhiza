# testscripts-refuse-work-in-repo

capability · `01M77W0Z9YDVWCQT2RSXGD8W91` · parent [self-check](../self-check.md)

## acceptance

A test in cmd/rhiza (go test -trimpath ./cmd/rhiza -run for it) calls the Setup check on a $WORK under a git work tree and under a .git dir and gets the error naming the path and GOTMPDIR, and on a $WORK outside any repo gets none; go test -trimpath ./cmd/rhiza -run TestScripts/ci-install passes

## statement

cmd/rhiza's testscript Setup fails a script before it runs when its $WORK resolves inside a git work tree or a git dir, with an error naming the path and saying to set GOTMPDIR outside the repository, so no script can write into a real repo

## Links

### Out

- decided-by [decision-run-tmp-outside-repo](../../../../decision/decision-run-tmp-outside-repo.md)

### In

- verifies [testscripts-refuse-work-in-repo-verified](../../../../evidence/testscripts-refuse-work-in-repo-verified.md)
- verifies [testscripts-refuse-work-in-repo-qa-1](../../../../evidence/testscripts-refuse-work-in-repo-qa-1.md)
- verifies [testscripts-refuse-work-in-repo-qa-2](../../../../evidence/testscripts-refuse-work-in-repo-qa-2.md)
- verifies [pr-488-fix-testscripts-run-inside-a](../../../../evidence/pr-488-fix-testscripts-run-inside-a.md)

## History

- `01M77VT6TMYCDVMX61ZVA8CS4C` Fix Testscripts run inside a ship run write into the real repo's main tree: the run's GOTMPDIR lies under ~/code/rhiza/.git (.git/rhiza/checkouts/...), so a script's $WORK sits inside the real repo, and git/rhiza commands in scripts that don't init their own repo resolve to ~/code/rhiza. Evidence: run 01M77DWN (cs 01M763Z4, PR #480), builder session 91d8e5b2 ran 'go test ./cmd/rhiza -run TestScripts/' at 2026-11-12T09:14:48Z; at 09:15:17 an untracked .github/workflows/rhiza.yml ('Written by rhiza ci install', from cmd/rhiza/testdata/script/ci-install.txtar, internal/api/ci.go:23 workflowPath) appeared in ~/code/rhiza's main tree. Moved aside to ~/orchestrator/stray-logs/2026-11-12/rhiza.yml by the orchestrator. Untracked files on main block close-out pulls and break the 'nothing uncommitted on main' rule. The agent later worked around it with GOTMPDIR=scratchpad. Fix ideas: runs set GOTMPDIR/TMPDIR outside any git repo (e.g. under the session scratchpad or ~/.cache/rhiza/tmp); TestMain refuses to start when $WORK resolves inside a git work tree other than its own.
