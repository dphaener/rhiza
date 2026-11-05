# machine-trust-unwritten-verified

evidence · `01M6MHE0RHZ6QK52Q518TD8TRA` · actor [claude-code-f8caac9c](../actor/claude-code-f8caac9c.md)

## result

PASS

## statement

At 9ce21088, verifier ran go test -trimpath -count=1 ./internal/machine/... -v: TestMachineNeverWritesFiles, TestMachineExportsOnlyReaders and TestNothingElseTouchesMachineFile (go/parser over non-test source; os/ioutil writers incl. WriteFile/Create/OpenFile without O_RDONLY/Mkdir*/Rename/Remove*/Chmod/Symlink/Link/Truncate; exports limited to Config, Daily, Dedicated, Load, Path, Read; no machine.toml literal passed to os/filepath/path calls outside internal/machine, aliased imports included) all PASS. go test -trimpath -count=1 ./cmd/rhiza -run TestScripts/machine-trust-unwritten PASS: with XDG_CONFIG_HOME set it runs machine show, api call machine.show, health and the refused session start --permission-mode bypassPermissions, first with no config (file and dir stay absent) then with trust = "daily" (cmp byte-identical against the kept copy, trust still daily). Mutation check in a scratch git-archive copy: adding os.WriteFile and an exported Write to internal/machine and an aliased filepath.Join(h, "rhiza", "machine.toml") in another package made all three guard tests fail with the expected messages.

## Links

### Out

- verifies [machine-trust-unwritten](../capability/rhiza/runtime/machine-trust-unwritten.md)

## History

- `01M6JYZ5Y4E56NMMC5HVKGC9J1` Constraint: Rhiza never writes or raises a machine's trust; the machine config is changed only by hand, so no agent, run or daemon can lift the bypass refusal through rhiza; record it as a constraint with a check that no op writes machine.toml (replaces draft change set 01M6GXMVKD6C14, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
