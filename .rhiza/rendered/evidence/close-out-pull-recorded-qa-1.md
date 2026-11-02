# close-out-pull-recorded-qa-1

evidence · `01M6DG4ATP35RJDNSFECAGR470` · actor [claude-code-6e06158d](../actor/claude-code-6e06158d.md)

## result

PASS

## statement

At ee31e06b, QA ran the real rhiza ship close-out by hand in a scratch repo with a bare origin, a change set and a ship state recording a merge commit. With README dirty, close-out exited 1 with 'main was not updated: ... has uncommitted changes' and wrote .git/rhiza/main-stale.json. rhiza health showed main warn '...; the last close-out could not pull' with a detail naming the change set, the time and the failure. rhiza inbox listed one main-stale item, and inbox open on it said there is no tab to open, next rhiza health. A second failing close-out kept the same since. After cleaning README and pushing a new commit to origin, close-out printed 'Updated main', pulled the commit and removed main-stale.json. Health main was back to ok and the inbox said 'Nothing waits on you'.

## Links

### Out

- verifies [close-out-pull-recorded](../capability/rhiza/orchestration/daemon/daemon-status/health-command/close-out-pull-recorded.md)

## History

- `01M6CRZSRFQ94MAJ8CRTEE44P2` rhiza health and a health app show the whole system's state on one screen
