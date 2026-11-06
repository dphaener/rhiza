# pr-327-fix-the-web-ui-s-overview-page

evidence · `01M6QXC20M1Q4SBFBB84K7MMEA` · actor [dphaener](../actor/dphaener.md)

## head

74c084b9547dbfa04ad65e23c35073e440165c43

## pr

https://github.com/dphaener/rhiza/pull/327

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #327 (https://github.com/dphaener/rhiza/pull/327) passed CI at 74c084b

## status

verified

## Links

### Out

- verifies [web-overview-fast](../capability/rhiza/web-monitor/web-overview-fast.md)

## History

- `01M6QSJHF5WH4SZKCTTYD0W0SS` Fix The web UI's Overview page takes about 5 s to load (owner: "it's incredibly slow"). Measured on the VM: / 5.2 s and /fragments/daemon 4.7 s, while the inbox, live-runs and milestones sections and /runs and /roadmap answer in 0.01-0.11 s. The daemon panel's source (WebDaemonSource, internal/api/web_daemon.go) calls runHealth on every render, and health runs the full model check (rhiza health alone takes 5.0 s here). The page and the daemon section's live refreshes must answer fast (well under 0.5 s) without running the check on each request, e.g. by reusing the daemon's last health result instead of recomputing it.
