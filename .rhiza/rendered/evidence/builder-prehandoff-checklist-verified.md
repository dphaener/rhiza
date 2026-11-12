# builder-prehandoff-checklist-verified

evidence · `01M7882BT9YWBEPB5A21QHPNC1` · actor [claude-code-2be75af1](../actor/claude-code-2be75af1.md)

## result

PASS

## statement

At 120eafed3: go test -trimpath -count=1 ./internal/guide -run Build passed (TestBuildPreHandoffChecklist asserts build.md names the acceptance commands, `rhiza ship generate`, `rhiza render`, the `Checks:` block and exit codes, and that Checks: sits above Change: and outside the final trailer paragraph; TestBuildTrailersInOneFinalParagraph still passes). go test -trimpath -count=1 ./internal/workflow passed, with TestBuildBuiltin asserting the builder prompt names the `pre-handoff checklist` and the `Checks:` block. make self-check in a clean detached worktree at the commit printed ok. Read the diff: build.md's builder section tells the builder to run the acceptance commands, rhiza ship generate and rhiza render, fold changed files in, and list each command with its exit code in a Checks: block above the trailers; the commit itself carries such a block.

## Links

### Out

- verifies [builder-prehandoff-checklist](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/builder-prehandoff-checklist.md)

## History

- `01M76XQJYPYK52PJRNAKCPPBJJ` Builders and fixers run the leaf's acceptance commands and the project's generate and render before their one commit, fold what they change into it, and report each command with its exit code in the commit; verifiers re-run what was reported. Evidence: 29 of 160 verifier blocker/major findings were generated, golden, format or self-check misses and 10 were commit, trailer or actor shape (prompt-gap rec 1).
