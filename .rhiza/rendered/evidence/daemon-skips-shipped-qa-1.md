# daemon-skips-shipped-qa-1

evidence · `01M69PWAY0MYJZ7TCYA24677J5` · actor [claude-code-6e84912c](../actor/claude-code-6e84912c.md)

## result

PASS

## statement

QA at 186386b0 with a scratch build (go build -o $TMPDIR/rhiza-qa ./cmd/rhiza) in a fresh git repo, fake rex on PATH, run env unset. Three accepted milestone change sets with leaves and hand-written journals: a succeeded ship-change run ending after the accept, a succeeded plan run ending after the accept, a succeeded ship-change run ending (2020) before the accept. rhiza daemon --once --max-runs=5 started runs of the plan and stale ones only (two rex window new calls) and skipped the shipped one. In a second repo on a feature branch: accepted, a succeeded ship-change journal ending after the accept, daemon --once skipped it as shipped; then rhiza land, commit, rhiza reopen --change <cs> --base main (landed -> accepted), and the next daemon --once started a run of it.

## Links

### Out

- verifies [daemon-skips-shipped](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-skips-shipped.md)

## History

- `01M69DFHXFG7197RGNYYPWS4A5` The daemon does not restart a change set whose latest run is a succeeded ship-change run that ended after the change set was last accepted or reopened
