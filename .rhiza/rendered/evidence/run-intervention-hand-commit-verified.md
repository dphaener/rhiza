# run-intervention-hand-commit-verified

evidence · `01M7R6WBAPKW7SQF882TR0JH0V` · actor [claude-code-97880160](../actor/claude-code-97880160.md)

## result

PASS

## statement

Verified builder commit a61a2b1d0 (capped scopes, RHIZA_/SHIP_/REX_ unset). Re-ran Checks: go test -trimpath -count=1 -p 2 ./internal/run -run 'HandCommit|Intervention' exit 0; ./internal/run exit 0; ./internal/api -run HandCommit exit 0 (TestHandCommitsByTrailer, TestHandCommitsOnlyRunBranch, TestHandCommitJournalledOnce pass); full ./internal/api exit 0 (214s, covers the Run|Watchdog|Daemon|Leaf|Review|Intervention subset); ./cmd/rhiza TestScripts/run-intervention-hand-commit exit 0; 22 run-intervention*/ship-change*/run-start*/run-resume* scripts exit 0 (stand-in for the 251-script pattern); go vet ./internal/api ./internal/run exit 0; rhiza ship generate and rhiza render on a detached checkout of a61a2b1d0 exit 0 with a clean tree. Acceptance met: a branch with a Leaf: commit of the change set, a Review: <run> commit, an Actor: commit of a run session and a merge of main yields 0 hand commits; a plain hand commit yields one hand-commit intervention naming it (actor = author); rescans (enter two, before finish, a fresh running on a resume) add none; Review: naming another run and Actor: of another run's session count as hand commits. must_not met: scanHandCommits only logs on any error or invalid intervention and never alters the run's result (TestHandCommitScanErrorChangesNothing); dedupe reads the journal from the store on every scan, so a restarted driver never journals one commit twice; git log is --first-parent Base..refs/heads/<run branch>, so only the run's own branch is scanned and main's commits brought in by merges are not; merges whose second parent is on origin/HEAD|origin/main|main are the run's (ship.rebase merges, never rebases). Daemon starts go through run start --foreground, so Base is recorded for them too. Note (no change asked): processActors also counts the engine process's own actor (falls back to git user.email's actor) as the run's, so a hand commit carrying the owner's own Actor: trailer in an owner-driven run is not journalled; this goes beyond the Decision's session-actor wording, and the commit message says why: ship step commits carry the process actor.

## Links

### Out

- verifies [run-intervention-hand-commit](../capability/rhiza/workflows/workflow-runs/run-intervention-hand-commit.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted
