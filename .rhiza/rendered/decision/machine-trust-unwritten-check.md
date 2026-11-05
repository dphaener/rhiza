# machine-trust-unwritten-check

decision · `01M6JZ5FER8SFM76M7DQWEJB92`

## reasoning

The question asks which check proves no op writes machine.toml. In the code, internal/machine (Path, Read, Load) is the only place that builds the machine config path, and its three callers (api/machine.go machine.show, api/health.go, api/session.go bypassCheck) only pass that path to machine.Read or os.Stat; session.go names $XDG_CONFIG_HOME/rhiza/machine.toml only as a message placeholder. machine-show's acceptance already checks describe --json lists no writing operation, but that only reads op summaries, not behaviour. A behavioural test alone cannot cover every op and every future one, so the structural source scan is what makes 'never' hold as code changes: it pins the file to one read-only package and fails the build the day something writes it or builds the path elsewhere. The testscript adds the observable proof the constraint's intent names, that the ops which read trust, including the bypass refusal an agent would try to lift, leave the file untouched and never create it. Both are hermetic (own XDG_CONFIG_HOME, no network, no Rex) and fit the repo's existing Go-test and testscript patterns (machine-trust-bypass.txtar, machine-show.txtar).

## statement

The leaf under runtime (beside machine-show) is proved by two checks, both in make check. (1) A Go test parses the module's non-test Go source (go/parser) and fails if internal/machine calls any file-writing or file-changing function (os.WriteFile, os.Create, os.OpenFile with write flags, os.Mkdir/MkdirAll, os.Rename, os.Remove/RemoveAll, os.Chmod, os.Symlink/Link, os.Truncate), if internal/machine exports anything but Path, Read, Load, Config and the trust constants, or if any file outside internal/machine names the machine config file other than as a display string (no filepath.Join or open of "machine.toml"), so machine.Path is the only way to the file and its package can only read it. (2) A testscript sets XDG_CONFIG_HOME, writes machine.toml with trust = "daily" and keeps a copy, then runs the ops that read trust: rhiza machine show (and api call machine.show), rhiza health, and the refused rhiza session start --permission-mode bypassPermissions; after them the file is byte-identical to the copy (cmp) and trust is still daily; a second pass with no machine config runs the same ops and the file still does not exist. machine-show's describe --json check (no operation writes the machine config) stays as it is.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-e4fa4f07-1](../question/question-e4fa4f07-1.md)
- derived-from [constraint-plan-interview](../source/constraint-plan-interview.md)

### In

- decided-by [machine-trust-unwritten](../capability/rhiza/runtime/machine-trust-unwritten.md)

## History

- `01M6JYZ5Y4E56NMMC5HVKGC9J1` Constraint: Rhiza never writes or raises a machine's trust; the machine config is changed only by hand, so no agent, run or daemon can lift the bypass refusal through rhiza; record it as a constraint with a check that no op writes machine.toml (replaces draft change set 01M6GXMVKD6C14, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
