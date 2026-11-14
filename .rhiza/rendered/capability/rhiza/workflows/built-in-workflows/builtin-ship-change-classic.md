# builtin-ship-change-classic

capability · `01M79KPE1SEZAAMTA1A37YWDR0` · parent [built-in-workflows](../built-in-workflows.md)

## acceptance

go test -trimpath ./internal/workflow -run TestBuiltinClassicUnchanged shows internal/workflow/builtin/ship-change-classic.toml equal to a golden copy of today's ship-change.toml (internal/workflow/testdata/ship-change-pre-m36.toml) apart from the name line, and BuiltinNames listing ship-change-classic; rhiza workflow check ship-change-classic is clean. Touches internal/workflow/builtin/ship-change-classic.toml, internal/workflow/testdata/ship-change-pre-m36.toml, internal/workflow/workflow_test.go.

## statement

The builtin ship-change as it is before milestone-36 ships, unchanged apart from its name, as the builtin ship-change-classic, so a project can pin it by name or in .rhiza/dispatch.toml.

## Links

### Out

- decided-by [conductor-promote-builtin-now](../../../../decision/conductor-promote-builtin-now.md)

### In

- depends-on [classic-tests-pinned](classic-tests-pinned.md)
- verifies [builtin-ship-change-classic-verified](../../../../evidence/builtin-ship-change-classic-verified.md)
- verifies [builtin-ship-change-classic-qa-2](../../../../evidence/builtin-ship-change-classic-qa-2.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
