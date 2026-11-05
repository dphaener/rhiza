# pr-238-usage-governor

evidence · `01M6JSHR82C24JHPSS0E4282S3` · actor [dphaener](../actor/dphaener.md)

## head

38fa3c3675d29e6cce65ed0675ad31aa89b87d60

## pr

https://github.com/dphaener/rhiza/pull/238

## result

1 check passed: check SUCCESS

## statement

Pull request #238 (https://github.com/dphaener/rhiza/pull/238) passed CI at 38fa3c3

## status

verified

## Links

### Out

- verifies [usage-governor](../capability/rhiza/orchestration/usage-governor.md)
- verifies [usage-window-tokens](../capability/rhiza/orchestration/usage-governor/usage-window-tokens.md)
- verifies [usage-windows](../capability/rhiza/orchestration/usage-governor/usage-windows.md)
- verifies [usage-settings](../capability/rhiza/orchestration/usage-governor/usage-settings.md)
- verifies [usage-learned-limit](../capability/rhiza/orchestration/usage-governor/usage-learned-limit.md)
- verifies [usage-hold](../capability/rhiza/orchestration/usage-governor/usage-hold.md)
- verifies [usage-hold-reported](../capability/rhiza/orchestration/usage-governor/usage-hold-reported.md)

## History

- `01M6HHMTEDFK2AWQKX50EB77DD` Usage governor: the daemon paces itself to a usage reserve. It reads usage per window from the Claude Code adapter's token records, holds new picks (with a skip reason in daemon status, health and inbox) when projected use would eat into a configured reserve (a daemon setting), and resumes when the window resets (internal/api/daemon.go, internal/adapter/claudecode). Carve the existing leaf capability usage-governor out of umbrella 01M5SWTC1923N3 (Area orchestration) into this change set rather than adding a duplicate
