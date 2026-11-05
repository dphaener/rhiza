# machine-trust-unwritten

capability · `01M6JZ8TSXBYMDF1G4PBXH3MRB` · parent [runtime](../runtime.md)

## acceptance

(1) A Go test parses the module's non-test Go source (go/parser) and fails if internal/machine calls any file-writing or file-changing function (os.WriteFile, os.Create, os.OpenFile with write flags, os.Mkdir/MkdirAll, os.Rename, os.Remove/RemoveAll, os.Chmod, os.Symlink/Link, os.Truncate), if internal/machine exports anything but Path, Read, Load, Config and the trust constants, or if any file outside internal/machine names the machine config file other than as a display string (no filepath.Join or open of "machine.toml"). (2) A testscript sets XDG_CONFIG_HOME, writes machine.toml with trust = "daily" and keeps a copy, runs rhiza machine show, rhiza api call machine.show, rhiza health and the refused rhiza session start --permission-mode bypassPermissions, then cmp shows the file byte-identical and trust still daily; a second pass with no machine config runs the same ops and the file still does not exist.

## statement

make check proves no rhiza code path writes the machine config: the file is reached only through internal/machine, which can only read it, and the ops that read trust leave it unchanged

## Links

### Out

- decided-by [machine-trust-unwritten-check](../../../decision/machine-trust-unwritten-check.md)
- decided-by [trust-hand-edited](../../../decision/trust-hand-edited.md)
- derived-from [constraint-plan-interview](../../../source/constraint-plan-interview.md)

### In

- constrains [machine-trust-never-written](../../../constraint/machine-trust-never-written.md)
- verifies [machine-trust-unwritten-verified](../../../evidence/machine-trust-unwritten-verified.md)
- verifies [machine-trust-unwritten-qa-1](../../../evidence/machine-trust-unwritten-qa-1.md)
- verifies [pr-268-constraint](../../../evidence/pr-268-constraint.md)

## History

- `01M6JYZ5Y4E56NMMC5HVKGC9J1` Constraint: Rhiza never writes or raises a machine's trust; the machine config is changed only by hand, so no agent, run or daemon can lift the bypass refusal through rhiza; record it as a constraint with a check that no op writes machine.toml (replaces draft change set 01M6GXMVKD6C14, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
