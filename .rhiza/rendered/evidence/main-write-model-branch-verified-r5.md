# main-write-model-branch-verified-r5

evidence · `01M76YAY2E5FBBQ2KY73XNP5WX` · actor [claude-code-f60ef53a](../actor/claude-code-f60ef53a.md)

## result

PASS

## statement

At c00162a77, in a detached scratch worktree outside ~/code (inside ~/code the live root daemon rewrites .rhiza-root/usage.live.json and trips the state-folder guard), ran go test -trimpath -count=1 ./internal/api -run 'TestModelBranch|TestHealthMain' -v: ok, all subtests PASS. With a local bare origin and a fake gh: AcceptShipsBranch shows an accept yields one pushed model/<full cs>-<entry> branch whose diff from origin/main is exactly one .rhiza/log entry, one PR create, the entry also in the main checkout, no worktree left, and a second write extends the same branch with no second PR; RecordApplyCarriesRender shows a --record apply's branch carries .rhiza/rendered/ and rhiza check on it reports no stale render; DryRunPushesNothing shows no branch and no gh call. Per decision-model-branch-fallback, NoOriginStaysLocal now runs an accept with no origin: entry written locally, no gh call, no stderr, no model/* branch, healthMain OK with no next step; FailingGHWarns runs an accept with a failing fake gh: exit 0, warning printed, no worktree left, healthMain bad naming that accept entry as an orphan, and the diff deletes the pushed branch when pr create fails so the entry is not hidden as pending. The leaf statement now names the branch with full change set and entry IDs, matching model_branch.go. Read the diff of c00162a77.

## Links

### Out

- verifies [main-write-model-branch](../capability/rhiza/orchestration/shipping/main-write-model-branch.md)

## History

- `01M76F267WMEGJDQFHHNM723TN` Nothing Rhiza writes is ever left uncommitted in a project's main checkout: every model write gets its own branch. Owner, 2026-11-11: "I don't feel like we should ever have anything uncommitted on main like this. Everything should get it's own branch." Evidence in ~/code/rhiza on 2026-11-11: untracked .rhiza/log/01M75W99SF6Q1T4QJM5GKA4PM1/ and 01M765854D7DGJDBW5CYGJ4SR1/ (the open entries of milestone requests 01M75W4E and 01M75YW5): newPlanChange (internal/api/run.go:855, openCtxWorkspace and w.openChangeSet at 882-894) opens a plan request's change set in the main checkout before the run's worktree exists, and the run then carries a copy on its branch, so main holds an untracked duplicate until close-out's pull (which git refuses when an untracked file would be overwritten). Also untracked: 01M75R9CDF2TVR2CSQM99SS2P3 (a request's draft withdrawn after its run was abandoned: the open entry was on main only, so the withdraw entry was too) and 01M6T24HZ2DATDFXP5EZBB0H0E (a change set 'list' opened and withdrawn on main on 2026-11-07), both never committed anywhere. Wanted: (1) run start --new (request plan runs, milestone records) creates the run's branch and worktree first and writes the open entry only there; (2) lifecycle and model writes run against the main checkout (rhiza change new, withdraw, accept/propose/reject from the CLI and the Inbox and web apps, apply, link, answers at a gate) never leave entries in the main checkout: they go to the change set's run worktree when one exists, else to a short-lived branch that the daemon commits, pushes and merges (as a model-only PR); (3) rhiza health and rhiza check report untracked .rhiza/log entries in a main checkout as a problem; (4) a migration or command that moves existing untracked main-checkout entries onto a branch. Related: 01M6RXN3BV80CHAQ1JCY9R7DTM (Inbox accept never reached the run), model-only CI (01M764CJ).
