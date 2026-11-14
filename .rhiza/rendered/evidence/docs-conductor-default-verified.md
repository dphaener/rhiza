# docs-conductor-default-verified

evidence · `01M7C4HDA5PSB8W5EDG8TT6AH3` · actor [claude-code-b4a673f9](../actor/claude-code-b4a673f9.md)

## result

PASS

## statement

At 7c35725d3 re-ran the Checks: block: go test -trimpath ./internal/guide -run TestWorkflowsText exit 0 (also with -count=1), go test -trimpath ./internal/doclinks exit 0. Read the diff: docs/runs-and-daemon.md now describes ship-change as the conductor default with its 10 steps (conduct, gate, signoff, push, wait-ci, merge, main-ci, close-out, ship, deploy-check), hand-back to conduct capped by handback_cap, and ship-change-classic as the build-then-review flow pinned via [ready] in .rhiza/dispatch.toml; the 'is build and then review' line is gone (grep finds none). docs/custom-workflows.md and internal/guide/text/workflows.md name ship-change-classic and the conductor default; workflows_test.go asserts ship-change-classic, conductor default and the hand-back phrase. Step names and hand-back wiring match internal/workflow/builtin/ship-change.toml. Only the four listed files plus the builder's actor record and its rendered files changed; no TODOs or stubs.

## Links

### Out

- verifies [docs-conductor-default](../capability/rhiza/workflows/built-in-workflows/docs-conductor-default.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
