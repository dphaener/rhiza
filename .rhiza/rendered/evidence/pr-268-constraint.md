# pr-268-constraint

evidence · `01M6MKT4DGNJ6N1FYE63XF4D3B` · actor [dphaener](../actor/dphaener.md)

## head

43e8794ea9ceef1504e4f49acd946dcd841a16cd

## pr

https://github.com/dphaener/rhiza/pull/268

## result

1 check passed: check SUCCESS

## statement

Pull request #268 (https://github.com/dphaener/rhiza/pull/268) passed CI at 43e8794

## status

verified

## Links

### Out

- verifies [machine-trust-unwritten](../capability/rhiza/runtime/machine-trust-unwritten.md)

## History

- `01M6JYZ5Y4E56NMMC5HVKGC9J1` Constraint: Rhiza never writes or raises a machine's trust; the machine config is changed only by hand, so no agent, run or daemon can lift the bypass refusal through rhiza; record it as a constraint with a check that no op writes machine.toml (replaces draft change set 01M6GXMVKD6C14, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
