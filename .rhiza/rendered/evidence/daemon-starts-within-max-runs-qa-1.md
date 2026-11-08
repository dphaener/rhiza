# daemon-starts-within-max-runs-qa-1

evidence · `01M6YEGH05PFF48RTTQJMR1NA9` · actor [claude-code-53686c59](../actor/claude-code-53686c59.md)

## result

PASS

## statement

QA at 3686b6ce. I built the binary with go build -trimpath into the session scratchpad and drove it by hand in a scratch root (HOME, git config, fake rex, RHIZA_TEST_DETACH_LOG so no real run starts). The setup had one live run (a process on another host), three ready change sets and a queued feature request and milestone request. At --max-runs 3, daemon --once started both requests' plan runs (run start --workflow plan-milestone ... --milestone, and run start --workflow plan) and no ship run: 1 live + 2 = 3. Two more --once polls at max 3 started nothing; detach.log stayed at 2 starts. daemon status showed 'runs: 1/3 live, 2 starting' and all three change sets 'waiting for a slot'. After the live run was removed in a second scratch root, --max-runs 3 with a resume still pending started exactly two ship runs (0 live + 1 resume pending + 2 = 3).

## Links

### Out

- verifies [daemon-starts-within-max-runs](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-starts-within-max-runs.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
