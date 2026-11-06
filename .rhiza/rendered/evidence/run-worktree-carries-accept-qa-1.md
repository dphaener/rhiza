# run-worktree-carries-accept-qa-1

evidence · `01M6S67PR1TKHB502BJ74ZC97A` · actor [claude-code-9bc184a1](../actor/claude-code-9bc184a1.md)

## result

PASS

## statement

QA at 90090bb2, hands-on with the built binary (go build -trimpath) in a scratch repo with a bare origin and a fake gh only for pr view, RHIZA_* env unset. Committed a proposed change set, ran rhiza accept in the main clone (entry left untracked), then rhiza run start --foreground: the new worktree showed State: accepted with the accept entry untracked; rhiza ship check passed and rhiza ship land committed and pushed, leaving a clean tree, and the accept entry file was in the pushed branch's tree on origin. For a second change set, a first run set up the worktree while it was proposed, then rhiza accept in the main clone and a second run on the existing worktree showed State: accepted and ship land landed it with the accept on the branch. A third run on the landed worktree copied nothing (clean porcelain, still landed). The same script with a binary built at the merge-base reproduced the bug: the worktree stayed proposed and ship land failed with 'cannot ship change set ...: it is proposed, not accepted'. go test -trimpath ./cmd/rhiza -run TestScripts/run-worktree-carries-accept passes.

## Links

### Out

- verifies [run-worktree-carries-accept](../capability/rhiza/orchestration/shipping/run-worktree-carries-accept.md)

## History

- `01M6RXNQFZX8YM80WFZJCXBRAJ` Fix Accepting a change set from the Inbox app never reaches the run that builds it. On 2026-11-06 the owner accepted Import procedure (01M6J36X) in the Inbox app at 14:31:10Z; the accept entry was written to the main clone's .rhiza/log but stayed uncommitted and unpushed. The daemon started the ship-change run 11 s later on a branch from origin, without the entry, and the run built, passed review and CI, then failed at land: 'cannot ship change set 01M6J36X...: it is proposed, not accepted' (the orchestrator copied the entry onto the run branch by hand to get it merged as #328). An accept from the Inbox app (and rhiza accept on the main clone generally) must reach the run's branch, e.g. by committing and pushing the accept to main or by the run carrying it, so the run never sees an accepted change set as proposed.
