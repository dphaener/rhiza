# builtin-ship-change-classic-verified

evidence · `01M7C0K7K3SJ21BZE1NTESJT3Q` · actor [claude-code-bef299ef](../actor/claude-code-bef299ef.md)

## result

PASS

## statement

At 424a3641d: go test -trimpath ./internal/workflow -run TestBuiltinClassicUnchanged exit 0 (also the whole internal/workflow package passes); a rhiza built from this commit runs `workflow check ship-change-classic`: valid (builtin:ship-change-classic), exit 0; rhiza ship generate and rhiza render exit 0 with no tracked changes; testscripts workflow-ls, workflow-check*, workflow-describe* and describe-golden pass. diff shows ship-change-classic.toml equal to internal/workflow/testdata/ship-change-pre-m36.toml except line 1 (name), and the golden equal byte for byte to the current builtin/ship-change.toml. The test checks both name lines, the rest of the file, BuiltinNames listing ship-change-classic and that it parses. The diff touches only the three named files plus this run's .rhiza log and rendered records, with no TODOs or stubs.

## Links

### Out

- verifies [builtin-ship-change-classic](../capability/rhiza/workflows/built-in-workflows/builtin-ship-change-classic.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
