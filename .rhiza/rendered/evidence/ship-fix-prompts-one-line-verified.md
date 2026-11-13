# ship-fix-prompts-one-line-verified

evidence · `01M79QSYP1GAE9TEDRRVGTKNAB` · actor [claude-code-69058fb1](../actor/claude-code-69058fb1.md)

## result

PASS

## statement

At ba8e6461c re-ran the Checks: block: go test -trimpath ./internal/workflow/... ok (TestShipFixPromptsOneLine asserts all six builtin ship-change agent steps' exact one-line prompts naming ci-fix, gate-fix, conflict-fix topics with no session done, Checks:, Review:, trailer or newline, and that .rhiza/workflows/ship-change.toml and docs.toml fix steps equal the builtin's); go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship-ci-fix|run-ship-ci-lint-fix|run-ship-conflict-fix|run-ship-change)' ok; make self-check exit 0; rhiza ship generate 0 changed files; rhiza render left no diff. rhiza guide ci-fix, gate-fix and conflict-fix all exist. Diff touches only the three workflow TOMLs, workflow_test.go and rhiza log/rendered actor records; no stubs or TODOs.

## Links

### Out

- verifies [ship-fix-prompts-one-line](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/ship-fix-prompts-one-line.md)

## History

- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
