# question-e4fa4f07-1

question · `01M6JZ2SJ8GNDPS2TA38ERNZG6`

## session

e4fa4f07-184f-4fd1-9134-01a1a5efce4d

## statement

The constraint 'Rhiza never writes or raises a machine's trust' needs a buildable check (only capabilities are build leaves, so I add one leaf under runtime, beside machine-show). Which check proves no op writes machine.toml?

## Links

### In

- answers [machine-trust-unwritten-check](../decision/machine-trust-unwritten-check.md)

## History

- `01M6JYZ5Y4E56NMMC5HVKGC9J1` Constraint: Rhiza never writes or raises a machine's trust; the machine config is changed only by hand, so no agent, run or daemon can lift the bypass refusal through rhiza; record it as a constraint with a check that no op writes machine.toml (replaces draft change set 01M6GXMVKD6C14, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
