# review-01m76f247z4y9t4j86ryf70hnt-r3-f2

question · `01M76FZ16C96GGBGK60879HHW2`

## statement

Plan review round 3 left finding f2 (contradiction on check-warns-orphan-log): In a main checkout, every untracked entry that no model branch or live run holds counts as an orphan. In testscript repos and other repos with no origin, every entry written is exactly that, because no model branch can exist there. rhiza check would then warn in nearly every existing script, and health-orphan-log-entries would mark the main part bad in every one, changing many goldens and contradicting the existing scripts' expected output. Neither the leaves nor the Decision says how a repo with no origin is treated. The plan should state whether a repo without origin/main is out of scope for orphan reporting, and the acceptances should cover it. The critic confirmed it: In a repo with no origin no model branch can exist, so every entry rhiza writes there is untracked and unheld, i.e. an orphan under check-warns-orphan-log and health-orphan-log-entries. 35 git-init testscripts run rhiza check after uncommitted writes (e.g. change-show-canonical-text.txtar, check-evidence-actor.txtar) and health scripts such as health.txtar and health-projects.txtar run in git-init repos, where the main part would turn bad (exit 1). Neither leaf nor the Decision says how a repo without origin/main is treated, and the acceptances do not cover it. How should it be settled?

## Links

### Out

- blocks [check-warns-orphan-log](../capability/rhiza/integrity/check/check-warns-orphan-log.md)

### In

- answers [orphan-log-needs-origin-main](../decision/orphan-log-needs-origin-main.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).
