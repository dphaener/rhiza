# constraint-plan-interview

source · `01M6JZ73TB639GE6J5WA410669`

## quotes

```json
[]
```

## raw_sha256

f9aaefb12ede23528f08242dbfcc11405b4d0e04d041e625ed02526bf3a771c0

## session

e4fa4f07-184f-4fd1-9134-01a1a5efce4d

## summary

Autonomous plan run for change set 01M6JYZ5 (Constraint). No owner present. The planner found internal/machine (Path, Read, Load) is the only code that builds the machine.toml path and its three callers (machine.show, health, the session bypass check) only read it; machine-show's acceptance checks describe --json lists no writing op. Only capabilities are build leaves, so the plan records a constraint 'Rhiza never writes or raises a machine's trust' plus one capability leaf under runtime that builds its check. The planner asked which check proves no op writes machine.toml (behavioural run-every-op, static package guard, or both). The arbiter decided machine-trust-unwritten-check: both a go/parser source scan (internal/machine calls no file-writing function, exports only Path/Read/Load/Config and trust constants, nothing outside it joins or opens machine.toml) and a testscript running machine show, api call machine.show, health and a refused bypass session start against a daily machine.toml, which stays byte-identical, and with no config, which stays absent. The change set serves milestone-25, relates to and replaces draft 01M6GXMVKD6C14, which is withdrawn.

## title

Constraint plan interview (machine trust never written)

## Links

### In

- derived-from [machine-trust-unwritten-check](../decision/machine-trust-unwritten-check.md)
- derived-from [machine-trust-never-written](../constraint/machine-trust-never-written.md)
- derived-from [machine-trust-unwritten](../capability/rhiza/runtime/machine-trust-unwritten.md)

## History

- `01M6JYZ5Y4E56NMMC5HVKGC9J1` Constraint: Rhiza never writes or raises a machine's trust; the machine config is changed only by hand, so no agent, run or daemon can lift the bypass refusal through rhiza; record it as a constraint with a check that no op writes machine.toml (replaces draft change set 01M6GXMVKD6C14, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
