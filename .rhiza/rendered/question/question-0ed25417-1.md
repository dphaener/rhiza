# question-0ed25417-1

question · `01M6RXVG3ECF82FS75G9HGYMPP`

## session

0ed25417-24bb-4e97-bfea-01a1b1db21df

## statement

How should an accept made in the main clone (Inbox app or rhiza accept) reach the ship-change run's branch? Today the run's worktree copies a change set's main-clone log entries only when the worktree lacks the change set entirely, so a committed proposed change set never gets its untracked accept entry. / This change set serves no intent yet. Which should it serve? Earlier daemon fixes (zombie drivers, owner-run gate accept, driver CPU) serve milestone-25 (backlog emptied by the daemon on its own); Inbox app fixes served m26-inbox-app.

## Links

### In

- answers [decision-run-carries-accept](../decision/decision-run-carries-accept.md)

## History

- `01M6RXNQFZX8YM80WFZJCXBRAJ` Fix Accepting a change set from the Inbox app never reaches the run that builds it. On 2026-11-06 the owner accepted Import procedure (01M6J36X) in the Inbox app at 14:31:10Z; the accept entry was written to the main clone's .rhiza/log but stayed uncommitted and unpushed. The daemon started the ship-change run 11 s later on a branch from origin, without the entry, and the run built, passed review and CI, then failed at land: 'cannot ship change set 01M6J36X...: it is proposed, not accepted' (the orchestrator copied the entry onto the run branch by hand to get it merged as #328). An accept from the Inbox app (and rhiza accept on the main clone generally) must reach the run's branch, e.g. by committing and pushing the accept to main or by the run carrying it, so the run never sees an accepted change set as proposed.
