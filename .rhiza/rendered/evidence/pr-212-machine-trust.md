# pr-212-machine-trust

evidence · `01M6H5VAZD5K1428FQ9PFCPKC4` · actor [dphaener](../actor/dphaener.md)

## head

deb4f8122d1e008be2be82e30db80885721588fd

## pr

https://github.com/dphaener/rhiza/pull/212

## result

1 check passed: check SUCCESS

## statement

Pull request #212 (https://github.com/dphaener/rhiza/pull/212) passed CI at deb4f81

## status

verified

## Links

### Out

- verifies [machine-config](../capability/rhiza/runtime/machine-config.md)
- verifies [bypass-on-dedicated](../capability/rhiza/runtime/bypass-on-dedicated.md)
- verifies [machine-show](../capability/rhiza/runtime/machine-show.md)
- verifies [health-machine-part](../capability/rhiza/runtime/health-machine-part.md)

## History

- `01M6FY9EXET6R8QN6GPEJTCF7B` Machine trust: a machine config can declare trust = dedicated (this VM does), and on such a machine the daemon and sessions may run with permission mode bypassPermissions; today bypassRefused (internal/api/session.go:437, session_run.go:153) refuses bypass everywhere because no machine config exists
