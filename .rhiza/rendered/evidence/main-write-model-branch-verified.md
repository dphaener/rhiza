# main-write-model-branch-verified

evidence · `01M76S70K1D1V9P4SFFT8336FC` · actor [claude-code-b3c56bbf](../actor/claude-code-b3c56bbf.md)

## result

PASS

## statement

Verifier at eb83f6d89: ran go test -trimpath -count=1 ./internal/api -run TestModelBranch -v; all five subtests PASS against a local bare origin and a fake gh. AcceptShipsBranch: accept pushes one model/<cs8>-<entry8> branch whose diff from origin/main is exactly one .rhiza/log entry, one pr create, temp worktree removed, entry also in the main checkout; a second write extends the same branch with no second PR. RecordApplyCarriesRender: a --record apply's branch carries .rhiza/rendered/ and rhiza check in a worktree of it reports no stale render. FailingGHWarns: write succeeds with a warning and no worktree left. DryRunPushesNothing: no branch, no gh call. NoOriginStaysLocal: no gh call. Adjacent ship/mirror tests also pass; go vet and gofmt clean. Read the diff (model_branch.go, workspace.go writeEntry hook, RHIZA_TEST_NOW gate wired through cmd/rhiza Getenv). The package run's trailing FAIL (.rhiza-root/usage.live.json changed) is the live root daemon rewriting that file every few seconds with no tests running, not this change.

## Links

### Out

- verifies [main-write-model-branch](../capability/rhiza/orchestration/shipping/main-write-model-branch.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).
