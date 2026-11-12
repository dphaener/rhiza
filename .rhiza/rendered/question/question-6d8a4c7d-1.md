# question-6d8a4c7d-1

question · `01M76F6QK59ZTSE4SHSV186QSC`

## session

6d8a4c7d-979b-46f4-96ad-01a1ccf15531

## statement

Change set 01M76F26 (nothing Rhiza writes is left uncommitted in the main checkout) serves no intent yet. Which should it serve? The earlier main-clone fixes it builds on (decision-run-carries-accept for the Inbox accept) served milestone-25, and milestone-25 is open on the roadmap. / A model write made in the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI, Inbox or web app, apply, link, gate answers) for a change set with no live run: where does its entry live until main has it? The daemon and the roadmap read change sets from the main checkout, so the main checkout must still see the write at once. (A change set with a live run gets the entry written into that run's worktree, and the run commits it as it does today.) / How are the untracked .rhiza/log entries already sitting in a main checkout (like 01M75R9C and 01M6T24H on ~/code/rhiza today) moved onto a branch? / How loudly do rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout?

## Links

### In

- answers [decision-main-writes-get-a-branch](../decision/decision-main-writes-get-a-branch.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).
