# machine-trust-never-written

constraint · `01M6JZ8F9AG1R89JSP0HGGR7TQ`

## acceptance

make check runs a Go source scan that pins machine.toml to internal/machine, which only reads it, and a testscript showing the ops that read trust leave machine.toml byte-identical and never create it

## statement

Rhiza never writes or raises a machine's trust: the machine config is changed only by hand, so no agent, run or daemon can lift the bypass refusal through rhiza.

## Links

### Out

- constrains [machine-show](../capability/rhiza/runtime/machine-show.md)
- constrains [machine-trust-unwritten](../capability/rhiza/runtime/machine-trust-unwritten.md)
- relates-to [bypass-only-on-dedicated](bypass-only-on-dedicated.md) — Never writing trust keeps the dedicated-only bypass rule from being lifted by rhiza itself
- derived-from [constraint-plan-interview](../source/constraint-plan-interview.md)

## History

- `01M6JYZ5Y4E56NMMC5HVKGC9J1` Constraint: Rhiza never writes or raises a machine's trust; the machine config is changed only by hand, so no agent, run or daemon can lift the bypass refusal through rhiza; record it as a constraint with a check that no op writes machine.toml (replaces draft change set 01M6GXMVKD6C14, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
