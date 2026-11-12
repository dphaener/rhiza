# daemon-merges-model-prs

capability · `01M76FAA6HBKA3Y4P11P3CZF4C` · parent [daemon](../daemon.md)

## acceptance

go test ./internal/api -run 'TestDaemonModelPR' with a fake gh and no run: a green model PR is squash-merged and its branch deleted, after which the main checkout is fast-forwarded to origin/main and holds no untracked copy of its entries; a red model PR is not merged and appears as an inbox item naming the change set; and go test ./internal/daemon -run 'TestSelfUpdate' : with an untracked .rhiza/log copy byte-identical to one origin/main adds, self-update drops it and fast-forwards

## statement

The daemon squash-merges each open model-only PR whose CI is green, deletes its branch, and then itself pulls the main checkout: it drops the main checkout's byte-identical untracked copies of the merged entries (dropIncomingLogCopies) and fast-forwards main, with no run close-out involved; the daemon's self-update also drops incoming copies before its fast-forward, so an untracked copy never blocks it; a model PR whose CI fails or conflicts is listed in the inbox

## Links

### Out

- decided-by [decision-main-writes-get-a-branch](../../../../decision/decision-main-writes-get-a-branch.md)

### In

- verifies [daemon-merges-model-prs-verified](../../../../evidence/daemon-merges-model-prs-verified.md)
- verifies [daemon-merges-model-prs-qa-1](../../../../evidence/daemon-merges-model-prs-qa-1.md)
- verifies [daemon-merges-model-prs-verified-5](../../../../evidence/daemon-merges-model-prs-verified-5.md)
- verifies [pr-464-nothing-rhiza-writes-is-ever](../../../../evidence/pr-464-nothing-rhiza-writes-is-ever.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).
