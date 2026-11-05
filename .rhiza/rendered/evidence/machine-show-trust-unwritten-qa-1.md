# machine-show-trust-unwritten-qa-1

evidence · `01M6MHNE2ZETT9XHFXCVA1BNBE` · actor [claude-code-5a9ef67a](../actor/claude-code-5a9ef67a.md)

## result

PASS

## statement

With the branch-head build, rhiza machine show printed 'trust: daily (no machine config; it would be <path>)' with no file and 'trust: daily (from <path>)' with a daily file; machine show --json and api call machine.show printed the same bytes ({file,path,trust}); describe --json names only the read op machine.show; neither run wrote the config.

## Links

### Out

- verifies [machine-show](../capability/rhiza/runtime/machine-show.md)

## History

- `01M6JYZ5Y4E56NMMC5HVKGC9J1` Constraint: Rhiza never writes or raises a machine's trust; the machine config is changed only by hand, so no agent, run or daemon can lift the bypass refusal through rhiza; record it as a constraint with a check that no op writes machine.toml (replaces draft change set 01M6GXMVKD6C14, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
