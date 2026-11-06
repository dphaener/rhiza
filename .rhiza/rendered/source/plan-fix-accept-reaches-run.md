# plan-fix-accept-reaches-run

source · `01M6RXYFVTNNZRZF7SS7TQD7PH`

## quotes

```json
[]
```

## raw_sha256

5ccd78e643c68537d5bcc01a8a7a32a379c28a8132846b04898f521af451d198

## session

0ed25417-24bb-4e97-bfea-01a1b1db21df

## summary

Autonomous plan run, no owner. The planner traced the bug: a run's worktree branches from the main clone's HEAD, and worktreeWorkspace copies the change set's main-clone log entries only when the worktree lacks the change set, so a committed proposed change set never gets its untracked accept entry and ship land fails with 'it is proposed, not accepted'. Two questions went to the arbiter: how the accept reaches the run (the run carries it, or accept pushes to main) and which intent the change set serves. The arbiter ruled the run carries it (worktree setup copies every entry of the change set the worktree lacks, committed on the run branch; accept never pushes to main) and that it serves milestone-25.

## title

Plan interview: Fix Accepting a change set (autonomous)

## Links

### In

- derived-from [decision-run-carries-accept](../decision/decision-run-carries-accept.md)

## History

- `01M6RXNQFZX8YM80WFZJCXBRAJ` Fix Accepting a change set from the Inbox app never reaches the run that builds it. On 2026-11-06 the owner accepted Import procedure (01M6J36X) in the Inbox app at 14:31:10Z; the accept entry was written to the main clone's .rhiza/log but stayed uncommitted and unpushed. The daemon started the ship-change run 11 s later on a branch from origin, without the entry, and the run built, passed review and CI, then failed at land: 'cannot ship change set 01M6J36X...: it is proposed, not accepted' (the orchestrator copied the entry onto the run branch by hand to get it merged as #328). An accept from the Inbox app (and rhiza accept on the main clone generally) must reach the run's branch, e.g. by committing and pushing the accept to main or by the run carrying it, so the run never sees an accepted change set as proposed.
