# main-write-model-branch

capability · `01M76FAA4YDVZV5S8GZB49W9M0` · parent [shipping](../shipping.md)

## acceptance

go test ./internal/api -run 'TestModelBranch' with a local bare origin and a fake gh: an accept in the main checkout yields a pushed model/<cs>-<entry> branch from origin/main holding exactly that entry and one PR request; a --record apply yields a branch holding its entry plus the re-rendered .rhiza/rendered/ files, on which rhiza check passes with no stale-render problem; a second write to the same change set adds a commit to the same branch and opens no second PR; the temporary worktree is removed; a dry run pushes nothing

## statement

A model write made in the main checkout for a change set with no live run still writes its entry into the main checkout's .rhiza/log and, in the same operation, commits it onto the short-lived branch model/<change set>-<entry>, named with the full change set and entry IDs, cut from origin/main in a temporary worktree, pushes it and opens a model-only PR; when the write changes the landed model (a --record apply or link, an add actor --record, an answer or accept that lands a record), the commit also carries rhiza render's output for .rhiza/rendered/ so the PR passes make check's self-check; a later write to the same change set while its model PR is open goes onto that same branch and PR

## Links

### Out

- decided-by [decision-main-writes-get-a-branch](../../../../decision/decision-main-writes-get-a-branch.md)

### In

- blocks [review-01m76f247z4y9t4j86ryf70hnt-r3-f1](../../../../question/review-01m76f247z4y9t4j86ryf70hnt-r3-f1.md)
- verifies [main-write-model-branch-verified](../../../../evidence/main-write-model-branch-verified.md)
- verifies [main-write-model-branch-qa-1](../../../../evidence/main-write-model-branch-qa-1.md)
- verifies [main-write-model-branch-verified-r2](../../../../evidence/main-write-model-branch-verified-r2.md)
- verifies [main-write-model-branch-qa-2](../../../../evidence/main-write-model-branch-qa-2.md)
- verifies [main-write-model-branch-verified-r5](../../../../evidence/main-write-model-branch-verified-r5.md)
- verifies [pr-464-nothing-rhiza-writes-is-ever](../../../../evidence/pr-464-nothing-rhiza-writes-is-ever.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).
