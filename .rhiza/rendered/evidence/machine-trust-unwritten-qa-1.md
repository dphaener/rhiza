# machine-trust-unwritten-qa-1

evidence · `01M6MHN610Z2J83DM2Z51DAD65` · actor [claude-code-5a9ef67a](../actor/claude-code-5a9ef67a.md)

## result

PASS

## statement

Built rhiza from the branch head and, in a scratch repo with its own HOME and XDG_CONFIG_HOME, ran machine show, api call machine.show, health and session start --permission-mode bypassPermissions (refused, naming machine.toml) with no machine config: XDG_CONFIG_HOME stayed empty (no rhiza dir, no file). Then with machine.toml trust = "daily": the same ops plus machine show --json left the file byte-identical (cmp), same mtime and mode, trust still daily. go test ./internal/machine and TestScripts/machine-trust-unwritten pass. Mutations: an aliased xos.WriteFile and an O_RDWR OpenFile in internal/machine plus a new export Raise fail TestMachineNeverWritesFiles and TestMachineExportsOnlyReaders; fp.Join(..., "machine.toml") in internal/api fails TestNothingElseTouchesMachineFile. A const-held "machine.toml" passed to filepath.Join outside internal/machine is not caught.

## Links

### Out

- verifies [machine-trust-unwritten](../capability/rhiza/runtime/machine-trust-unwritten.md)

## History

- `01M6JYZ5Y4E56NMMC5HVKGC9J1` Constraint: Rhiza never writes or raises a machine's trust; the machine config is changed only by hand, so no agent, run or daemon can lift the bypass refusal through rhiza; record it as a constraint with a check that no op writes machine.toml (replaces draft change set 01M6GXMVKD6C14, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
