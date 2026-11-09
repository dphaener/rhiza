# health-runs-against-cap-qa-1

evidence · `01M6YEGHK7CZFJC28E9T9A2ENT` · actor [claude-code-53686c59](../actor/claude-code-53686c59.md)

## result

PASS

## statement

QA at 3686b6ce, with the built binary running a long-lived daemon in a scratch root (one live run, two unjournaled plan-run starts). At --max-runs 3, rhiza health showed 'ok daemon running ...' with 'runs: 1/3 live, 2 starting', and health --json's daemon part carried live 1, starting 2, max 3, state ok. At --max-runs 1 it showed 'bad daemon running, ... 1 live runs; over rhiza.daemon.maxRuns 1', 'runs: 1/1 live, 2 starting' and a 'live run: <id>, ' line, health exited 1, and --json gave state bad with the same live, starting and max.

## Links

### Out

- verifies [health-runs-against-cap](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-runs-against-cap.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
