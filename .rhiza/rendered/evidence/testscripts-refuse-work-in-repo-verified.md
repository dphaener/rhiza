# testscripts-refuse-work-in-repo-verified

evidence · `01M77WFTZ4PB4SC11STRPKDQA6` · actor [claude-code-f3b52036](../actor/claude-code-f3b52036.md)

## result

PASS

## statement

At 561ecd015 I read the diff: setup() in cmd/rhiza/main_test.go now calls checkWorkOutsideRepo(env.WorkDir) before anything else; it resolves symlinks and walks ancestors, failing on a .git path element, an ancestor holding .git (work tree), or a bare-repo dir (HEAD+objects+refs), with errors naming the path and 'set GOTMPDIR outside the repository'. go test -trimpath ./cmd/rhiza -run 'TestCheckWorkOutsideRepo|TestScripts/ci-install' -v -count=1 with GOTMPDIR in my scratchpad: TestCheckWorkOutsideRepo/work_tree, /git_dir, /bare_repo (errors containing GOTMPDIR and the path) and /outside (nil) PASS, TestScripts/ci-install PASS. Run with the run's own GOTMPDIR under ~/code/rhiza/.git/rhiza/runs/.../tmp, ci-install fails before running with 'testscript $WORK $WORK is inside a git dir: set GOTMPDIR outside the repository', so the guard fires in the real case. Note: testscript's log rewrites the work dir to $WORK, so the printed path reads '$WORK' rather than the absolute path; the error itself carries the real path.

## Links

### Out

- verifies [testscripts-refuse-work-in-repo](../capability/rhiza/integrity/self-check/testscripts-refuse-work-in-repo.md)

## History

- `01M77VT6TMYCDVMX61ZVA8CS4C` Fix Testscripts run inside a ship run write into the real repo's main tree: the run's GOTMPDIR lies under ~/code/rhiza/.git (.git/rhiza/checkouts/...), so a script's $WORK sits inside the real repo, and git/rhiza commands in scripts that don't init their own repo resolve to ~/code/rhiza. Evidence: run 01M77DWN (cs 01M763Z4, PR #480), builder session 91d8e5b2 ran 'go test ./cmd/rhiza -run TestScripts/' at 2026-11-12T09:14:48Z; at 09:15:17 an untracked .github/workflows/rhiza.yml ('Written by rhiza ci install', from cmd/rhiza/testdata/script/ci-install.txtar, internal/api/ci.go:23 workflowPath) appeared in ~/code/rhiza's main tree. Moved aside to ~/orchestrator/stray-logs/2026-11-12/rhiza.yml by the orchestrator. Untracked files on main block close-out pulls and break the 'nothing uncommitted on main' rule. The agent later worked around it with GOTMPDIR=scratchpad. Fix ideas: runs set GOTMPDIR/TMPDIR outside any git repo (e.g. under the session scratchpad or ~/.cache/rhiza/tmp); TestMain refuses to start when $WORK resolves inside a git work tree other than its own.
