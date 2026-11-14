# docs-conductor-default

capability · `01M7AHXH9AZCZPEX26PHDJ1V5M` · parent [built-in-workflows](../built-in-workflows.md)

## acceptance

go test -trimpath ./internal/guide -run TestWorkflowsText shows the workflows topic naming ship-change-classic and the conductor default; go test -trimpath ./internal/doclinks passes; docs/runs-and-daemon.md no longer says ship-change is build and then review. Touches docs/runs-and-daemon.md, docs/custom-workflows.md, internal/guide/text/workflows.md, internal/guide/workflows_test.go.

## statement

docs/runs-and-daemon.md, docs/custom-workflows.md and the workflows guide topic describe ship-change as the conductor default (its 10 steps and the hand-back to conduct) and ship-change-classic as the old build-then-review flow a project can pin in .rhiza/dispatch.toml.

## Links

### Out

- depends-on [builtin-conductor-default](builtin-conductor-default.md)
- decided-by [conductor-promote-builtin-now](../../../../decision/conductor-promote-builtin-now.md)

### In

- verifies [docs-conductor-default-verified](../../../../evidence/docs-conductor-default-verified.md)
- verifies [docs-conductor-default-qa-2](../../../../evidence/docs-conductor-default-qa-2.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
