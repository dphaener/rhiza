# daemon-skip-reasons-qa-1

evidence · `01M6950ZH1PFCFXVJWHHBY6ZVE` · actor [claude-code-150282ee](../actor/claude-code-150282ee.md)

## result

PASS

## statement

At 9ea6aab built rhiza to a scratch path and, in a scratch git repo with RHIZA_* and REX_* unset and a logging fake rex on PATH, accepted Alpha and Beta (milestone, leaf), Orphan (no milestone), Leafless (no leaves) and Waiting (depends on Orphan). daemon --once --max-runs=1 started only Alpha (one rex window new) and .git/rhiza/daemon.state.json held skips: Leafless no-leaves, Waiting not-ready with detail 'depends on 01M694VFHJFPQ0 (accepted)', Orphan no-milestone, Beta no-slot, each with its title. After planting a running journal for Alpha (process pid 1 on this host) and a failed journal for Beta, a further daemon --once --max-runs=2 started nothing (rex log empty) and saved Alpha live and Beta failed beside the other three. held and unjournaled not exercised by hand (they need a start the same daemon process sees unjournaled).

## Links

### Out

- verifies [daemon-skip-reasons](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-skip-reasons.md)

## History

- `01M686RS47EK7F4H2WY15TNAQ5` rhiza daemon status names each accepted change set the daemon does not start and why
