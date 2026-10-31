# shipped-skip-reason-qa-1

evidence · `01M69PWCR1G6ZFRZPHB5D96HA2` · actor [claude-code-6e84912c](../actor/claude-code-6e84912c.md)

## result

PASS

## statement

QA at 186386b0 with a scratch build in a fresh repo, fake rex: after daemon --once, rhiza daemon status printed 'Skipped: <cs> shipped: its latest run shipped it' for the change set whose latest run is a succeeded ship-change run after its accept, and daemon status --json listed it with reason shipped and detail []; the plan-run and pre-accept ship-run change sets were not among the skips. After a reopen the same change set was no longer skipped and was started.

## Links

### Out

- verifies [shipped-skip-reason](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-skip-reasons/shipped-skip-reason.md)

## History

- `01M69DFHXFG7197RGNYYPWS4A5` The daemon does not restart a change set whose latest run is a succeeded ship-change run that ended after the change set was last accepted or reopened
