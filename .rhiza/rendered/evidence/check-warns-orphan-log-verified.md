# check-warns-orphan-log-verified

evidence · `01M76QEGKJNYQPEZBB4118K1NH` · actor [claude-code-36ad4e90](../actor/claude-code-36ad4e90.md)

## result

PASS

## statement

Verifier (attempt 2) at c55e23ef9: read the diff of 7a96ef5b4 (internal/api/check.go orphanLogEntries, internal/check/check.go Orphans warn finding with Next 'rhiza model commit-untracked') and ran go test -trimpath -count=1 ./cmd/rhiza -run TestScripts/check-orphan-log (ok) and go test -trimpath ./internal/check/... (ok); go build ./... and go vet on internal/api and internal/check clean. The script asserts: a git-init repo with no origin gives ok and no warning; with origin/main an untracked change-set entry in the main checkout gives exit 0, 'ok', a 'warn log .rhiza/log/<cs>/<entry>.json: untracked log entry that no model branch or live run holds' line naming the entry and 'next: rhiza model commit-untracked', and --json ok:true with a warn finding; the same entry copied into a linked worktree gives no warning in that worktree (and holds the main one); removing the worktree makes it an orphan again; a model/* branch holding the path also silences it. orphanLogEntries returns nil in linked worktrees (git-dir != common-dir) and without refs/remotes/origin/main, matching orphan-log-needs-origin-main.

## Links

### Out

- verifies [check-warns-orphan-log](../capability/rhiza/integrity/check/check-warns-orphan-log.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).
