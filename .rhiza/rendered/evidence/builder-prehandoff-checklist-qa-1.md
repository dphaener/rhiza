# builder-prehandoff-checklist-qa-1

evidence · `01M788Q85YT2ZX81MG94F21NCD` · actor [claude-code-22403536](../actor/claude-code-22403536.md)

## result

PASS

## statement

QA at 95fee8aca with a binary built from the head: `rhiza guide build` prints the builder checklist (acceptance commands, `rhiza ship generate`, `rhiza render`, fold changes in, a Checks: block with exit codes above the trailers) and `rhiza workflow describe build --json` shows the builder prompt naming the pre-handoff checklist. In a scratch git repo, `rhiza ship generate` exited 0 both unset (nothing to generate) and with rhiza.ship.generate set (staged the changed file), and `rhiza render` exited 0. A commit made from the guide's example message parses with `git interpret-trailers --parse` to exactly Change/Actor/Leaf/Co-Authored-By, so the Checks: block is not read as trailers; `rhiza check` on this branch claims its Checks:-carrying commits. go test -trimpath ./internal/guide ./internal/workflow pass.

## Links

### Out

- verifies [builder-prehandoff-checklist](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/builder-prehandoff-checklist.md)

## History

- `01M76XQJYPYK52PJRNAKCPPBJJ` Builders and fixers run the leaf's acceptance commands and the project's generate and render before their one commit, fold what they change into it, and report each command with its exit code in the commit; verifiers re-run what was reported. Evidence: 29 of 160 verifier blocker/major findings were generated, golden, format or self-check misses and 10 were commit, trailer or actor shape (prompt-gap rec 1).
