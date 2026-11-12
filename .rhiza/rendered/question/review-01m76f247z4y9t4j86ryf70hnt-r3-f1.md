# review-01m76f247z4y9t4j86ryf70hnt-r3-f1

question · `01M76FZ15DP1PVV3Y4R7KP5AQ4`

## statement

Plan review round 3 left finding f1 (contradiction on main-write-model-branch): Every model write in a main checkout with no live run must, in the same operation, cut a branch from origin/main, push it and open a PR, but neither the leaf nor decision-main-writes-get-a-branch says what happens when the repo has no origin remote, no origin/main or no gh. That covers every testscript repo (CLAUDE.md: no test touches the network or gh; each script's repo is a plain, non-linked checkout) and any local-only Rhiza repo. As written, every rhiza change new, accept or apply in those repos would fail, or would need network and gh, which contradicts the test constraints and breaks the existing scripts. The leaf should state the fallback (e.g. no origin or gh: write locally and skip the model branch, or leave the entry to health/commit-untracked), and its acceptance should cover it. The critic confirmed it: Both main-write-model-branch and decision-main-writes-get-a-branch (2) make the branch, push and PR part of the same operation for every main-checkout write with no live run, with no fallback. Of 359 testscripts, 270 git init a plain main checkout and only 51 add an origin (e.g. change-show-canonical-text.txtar: git init, rhiza change new, add, propose, accept, land, no origin); apply.txtar has no git at all. CLAUDE.md forbids tests touching the network or gh. As written these writes would fail or need gh; the leaf's acceptance only covers a repo with a bare origin and fake gh, so the no-origin/no-gh case is undecided. How should it be settled?

## Links

### Out

- blocks [main-write-model-branch](../capability/rhiza/orchestration/shipping/main-write-model-branch.md)

### In

- answers [decision-model-branch-fallback](../decision/decision-model-branch-fallback.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).
