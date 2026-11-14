# handback-target-survives-include-verified

evidence · `01M7C18GNT3XSYA7P47H05KHKC` · actor [claude-code-39725c4e](../actor/claude-code-39725c4e.md)

## result

PASS

## statement

Verified at 748636e58. Re-ran the Checks block: go test -trimpath ./internal/workflow -run TestIncludeHandbackTarget (PASS, exit 0), go test -trimpath ./internal/run -run TestHandbackIncluded (PASS, exit 0), rhiza ship generate (exit 0, 0 files changed), rhiza render (no diff). Diff read: TestIncludeHandbackTarget checks that an included gate's on.fail handback:conduct parses as handback:inner/conduct, that inner/conduct is a flat step, that handback:nosuch is refused with an error naming nosuch and "inner", and that done/failed still map to the including step's on.pass/on.fail. TestHandbackIncluded runs an including workflow whose included gate fails once and checks it hands back to inner/conduct (one Handback, target inner/conduct) and succeeds. The mapping was already in internal/workflow/include.go (inlined.outer cuts HandbackPrefix and returns handback:<workflow>/<step>), so the commit adds only tests, the builder's actor record and its rendered files; nothing outside the leaf changed, and there are no stubs or TODOs.

## Links

### Out

- verifies [handback-target-survives-include](../capability/rhiza/workflows/built-in-workflows/handback-target-survives-include.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
