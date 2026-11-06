# decision-run-carries-accept

decision · `01M6RXXPWWKZCVSTCWT9NB149S`

## reasoning

internal/api/session.go worktreeWorkspace returns early when the worktree has the change set at all, so a committed proposed change set never picks up the untracked accept entry; dropping that early return is a small fix in code that already exists, and copyChangeSetLog is idempotent because entries are immutable files named by ID. Carrying the entry on the run branch is what the orchestrator did by hand to merge #328, and ship_close.go dropIncomingLogCopies already removes untracked main-clone log files that origin/main then holds byte for byte, so the close-out pull still works. Committing and pushing to main from the Inbox app or rhiza accept would make an app write to origin's main outside a run's PR and CI, and it would dirty or race the main clone that close-out pulls into. For the intent: the failure stopped the daemon from shipping an accepted change set on its own, which is milestone-25's promise, like the earlier daemon fixes (owner-run gate accept, zombie drivers, driver CPU); the Inbox app only wrote the entry correctly, so m26-inbox-app does not fit.

## statement

The run carries the accept: when a run's worktree is set up for change set cs, worktreeWorkspace copies every entry of cs that the main clone's log has and the worktree's lacks, even when the worktree already has the change set (copyChangeSetLog already skips files that exist), and the run commits those entries on its branch so the land sees the change set as accepted. Accepting in the Inbox app or with rhiza accept does not commit or push to main. This change set serves milestone-25.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-0ed25417-1](../question/question-0ed25417-1.md)
- derived-from [plan-fix-accept-reaches-run](../source/plan-fix-accept-reaches-run.md)

### In

- decided-by [run-worktree-carries-accept](../capability/rhiza/orchestration/shipping/run-worktree-carries-accept.md)

## History

- `01M6RXNQFZX8YM80WFZJCXBRAJ` Fix Accepting a change set from the Inbox app never reaches the run that builds it. On 2026-11-06 the owner accepted Import procedure (01M6J36X) in the Inbox app at 14:31:10Z; the accept entry was written to the main clone's .rhiza/log but stayed uncommitted and unpushed. The daemon started the ship-change run 11 s later on a branch from origin, without the entry, and the run built, passed review and CI, then failed at land: 'cannot ship change set 01M6J36X...: it is proposed, not accepted' (the orchestrator copied the entry onto the run branch by hand to get it merged as #328). An accept from the Inbox app (and rhiza accept on the main clone generally) must reach the run's branch, e.g. by committing and pushing the accept to main or by the run carrying it, so the run never sees an accepted change set as proposed.
