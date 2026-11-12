# main-writes-branch-plan

source · `01M76F9M97Z9XNQND231WG7SDJ`

## quotes

```json
[
  {
    "by": "owner",
    "text": "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch.",
    "where": "2026-11-11, quoted in the change set's intent"
  }
]
```

## raw_sha256

bdfc79e2acd7f40930b1933b02f78d45b2599eac3e5f078410cc39dcb386d337

## session

6d8a4c7d-979b-46f4-96ad-01a1ccf15531

## summary

Autonomous plan run, so no owner interview. The planner read the change set's intent, the code (newPlanChange in internal/api/run.go opens a plan request's change set in the main checkout before the run's worktree exists; workspace.writeEntry in internal/api/workspace.go is the one place model entries are written), and the earlier Decisions on main-clone entries (decision-run-carries-accept, close-out-untracked-log-and-retry). It raised four choices to the arbiter: which intent the change set serves, where a main-checkout write lives until main has it, how existing untracked entries move onto a branch, and how loudly health and check report them. The arbiter answered with Decision decision-main-writes-get-a-branch: serve milestone-25; a main-checkout write still lands in the main checkout's log and in the same operation is committed onto a short-lived model branch from origin/main in a temporary worktree, pushed and opened as a model-only PR that the daemon squash-merges on green (a landed write's commit also carries rhiza render's output), later writes to the same change set joining that branch; the daemon (each tick) and a CLI command move orphaned untracked entries onto model branches the same way; health reports an orphan as a problem and lists entries held by a model branch as pending; rhiza check warns in a main checkout only. The plan's leaves: run start --new opens the change set in the run's worktree, where the run commits it, and writes the same bytes into the main checkout's log as a copy held by the live run; writes for a change set with a live run also go into that run's worktree; main-checkout writes get a model branch and PR; the daemon merges model PRs, then drops the main checkout's copies of their entries and fast-forwards main itself, as self-update also does; orphans are committed by the daemon and by rhiza model commit-untracked; health and check report orphans.

## title

Plan of change set 01M76F26: nothing Rhiza writes is left uncommitted in the main checkout

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).
