# run-worktree-carries-accept

capability · `01M6RXZ1WMDHH9H5ZH8T7G0VR2` · parent [shipping](../shipping.md)

## acceptance

A new testscript (go test -trimpath ./cmd/rhiza -run TestScripts/run-worktree-carries-accept) commits a proposed change set with leaves, accepts it in the main clone so the accept entry stays untracked, then starts a run on it (a new worktree, and again with the worktree already existing on its branch) and shows change show in the worktree reports accepted and the accept entry file is there; rhiza ship land in that worktree commits the accept entry on the branch. go test -trimpath ./internal/api/... passes.

## statement

When a run or session sets up its worktree for a change set, every log entry of that change set the main clone has and the worktree lacks is copied in, even when the worktree already has the change set, so an accept made in the main clone (Inbox app or rhiza accept, never pushed) reaches the run and ship land commits it on the run's branch.

## Links

### Out

- decided-by [decision-run-carries-accept](../../../../decision/decision-run-carries-accept.md)

### In

- verifies [run-worktree-carries-accept-verified](../../../../evidence/run-worktree-carries-accept-verified.md)
- verifies [run-worktree-carries-accept-qa-1](../../../../evidence/run-worktree-carries-accept-qa-1.md)
- verifies [pr-334-fix-accepting-a-change-set](../../../../evidence/pr-334-fix-accepting-a-change-set.md)

## History

- `01M6RXNQFZX8YM80WFZJCXBRAJ` Fix Accepting a change set from the Inbox app never reaches the run that builds it. On 2026-11-06 the owner accepted Import procedure (01M6J36X) in the Inbox app at 14:31:10Z; the accept entry was written to the main clone's .rhiza/log but stayed uncommitted and unpushed. The daemon started the ship-change run 11 s later on a branch from origin, without the entry, and the run built, passed review and CI, then failed at land: 'cannot ship change set 01M6J36X...: it is proposed, not accepted' (the orchestrator copied the entry onto the run branch by hand to get it merged as #328). An accept from the Inbox app (and rhiza accept on the main clone generally) must reach the run's branch, e.g. by committing and pushing the accept to main or by the run carrying it, so the run never sees an accepted change set as proposed.
