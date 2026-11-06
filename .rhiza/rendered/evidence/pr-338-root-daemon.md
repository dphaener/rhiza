# pr-338-root-daemon

evidence · `01M6SMGV5CZG5REAB9JSZK428M` · actor [dphaener](../actor/dphaener.md)

## head

c254a7ef0dfc725445ef0abde94e1363f235a5d4

## pr

https://github.com/dphaener/rhiza/pull/338

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #338 (https://github.com/dphaener/rhiza/pull/338) passed CI at c254a7e

## status

verified

## Links

### Out

- verifies [root-daemon-lock](../capability/rhiza/orchestration/daemon/root-daemon-lock.md)
- verifies [root-daemon-polls-projects](../capability/rhiza/orchestration/daemon/root-daemon-polls-projects.md)
- verifies [root-daemon-round-robin](../capability/rhiza/orchestration/daemon/root-daemon-round-robin.md)
- verifies [root-daemon-settings-scope](../capability/rhiza/orchestration/daemon/root-daemon-settings-scope.md)
- verifies [root-daemon-usage-governor](../capability/rhiza/orchestration/usage-governor/root-daemon-usage-governor.md)
- verifies [root-daemon-service](../capability/rhiza/orchestration/daemon/root-daemon-service.md)
- verifies [root-daemon-self-update](../capability/rhiza/orchestration/daemon/root-daemon-self-update.md)
- verifies [root-daemon-rex-session](../capability/rhiza/orchestration/daemon/root-daemon-rex-session.md)
- verifies [daemon-settings](../capability/rhiza/orchestration/daemon/daemon-settings.md)
- verifies [usage-settings](../capability/rhiza/orchestration/usage-governor/usage-settings.md)
- verifies [web-serve-host](../capability/rhiza/web-monitor/web-serve-loopback/web-serve-host.md)
- verifies [daemon-picks](../capability/rhiza/orchestration/daemon/daemon-picks.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
