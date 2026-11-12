# builder-prehandoff-checklist-verified-2

evidence · `01M7897DG8M2FF9R625R09F2TW` · actor [claude-code-e26b0d44](../actor/claude-code-e26b0d44.md)

## result

PASS

## statement

Verifier, attempt 2, at d0e9c974d: `go test -trimpath ./internal/guide -run Build -v` passes (TestBuildText, TestBuildPreHandoffChecklist, TestBuildCommitsOwnActorRecord, TestBuildTrailersInOneFinalParagraph, TestBuildVerifierRerunsChecks); `go test -trimpath ./internal/guide` and `go test -trimpath ./internal/workflow` exit 0, and workflow_test.go:626 asserts the build.toml builder prompt names the 'pre-handoff checklist' and the `Checks:` block. In build.md, the builder section names the acceptance commands, `rhiza ship generate`, `rhiza render`, folding in the files they change, and a `Checks:` block of commands with exit codes. The example puts Checks: above the one final trailer paragraph, which keeps Change/Actor/Leaf/Co-Authored-By together. The fix commit restores the three majors from attempt 1: 'build only that leaf', 'after its recording command succeeded' and '`Co-Authored-By:` included, in the one final paragraph, no blank line'. build_test.go now pins all three. In a clean detached worktree, `rhiza check` reports only stale-render failures for run actors. All but this commit's own actor were already failing at the parent a56cf9a13, and the run's post-build render covers them.

## Links

### Out

- verifies [builder-prehandoff-checklist](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/builder-prehandoff-checklist.md)

## History

- `01M76XQJYPYK52PJRNAKCPPBJJ` Builders and fixers run the leaf's acceptance commands and the project's generate and render before their one commit, fold what they change into it, and report each command with its exit code in the commit; verifiers re-run what was reported. Evidence: 29 of 160 verifier blocker/major findings were generated, golden, format or self-check misses and 10 were commit, trailer or actor shape (prompt-gap rec 1).
