# daemon-skips-shipped-verified

evidence · `01M69NHEM4M0A3KDGQ640DSM3Q` · actor [claude-code-0cb5428c](../actor/claude-code-0cb5428c.md)

## result

PASS

## statement

At 03e9795e, TestPickSkipsShipped in internal/daemon/pick_test.go has the three table cases the acceptance names: a succeeded ship-change run after the accept is skipped, a succeeded plan run is picked, and a succeeded ship-change run followed by a reopen is picked. go test ./internal/daemon -run 'TestPickSkipsShipped|TestSkips' -count=1 passed, and go test ./... passed for every package. With go test -overlay mutations of pick.go (the checkout unchanged), each case failed on its own: removing the shipped skip failed the skip case, ignoring Since failed the reopen case, and dropping the workflow check failed the plan case. make check's self-check failed only on stale .rhiza/rendered actor and index files from this run's actor records, not on code.

## Links

### Out

- verifies [daemon-skips-shipped](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-skips-shipped.md)

## History

- `01M69DFHXFG7197RGNYYPWS4A5` The daemon does not restart a change set whose latest run is a succeeded ship-change run that ended after the change set was last accepted or reopened
