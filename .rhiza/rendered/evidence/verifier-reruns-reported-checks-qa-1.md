# verifier-reruns-reported-checks-qa-1

evidence · `01M788Q89CAA8R2G8C4S8V53D8` · actor [claude-code-22403536](../actor/claude-code-22403536.md)

## result

PASS

## statement

QA at 95fee8aca with a binary built from the head: `rhiza guide build`'s Verifier section says to re-run every command in the commit's Checks: block and every acceptance command, and that a different exit code or no Checks: block is a finding naming the command. Re-running the Checks: commands of this branch's 5f8bfd112 in scope (go test -trimpath ./internal/guide ./internal/workflow) exited 0 as reported. go test -trimpath ./internal/guide passes.

## Links

### Out

- verifies [verifier-reruns-reported-checks](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/verifier-reruns-reported-checks.md)

## History

- `01M76XQJYPYK52PJRNAKCPPBJJ` Builders and fixers run the leaf's acceptance commands and the project's generate and render before their one commit, fold what they change into it, and report each command with its exit code in the commit; verifiers re-run what was reported. Evidence: 29 of 160 verifier blocker/major findings were generated, golden, format or self-check misses and 10 were commit, trailer or actor shape (prompt-gap rec 1).
