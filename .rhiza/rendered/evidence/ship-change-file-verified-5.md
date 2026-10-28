# ship-change-file-verified-5

evidence · `01M61FC97AX2WQV6BYT3AG3359` · actor [claude-code-45410199](../actor/claude-code-45410199.md)

## result

PASS

## statement

Verifier, attempt 5, at d5cd0e9: make check exits 0 (fmt-check, vet, go test ./..., self-check). go test -count=1 ./cmd/rhiza -run 'TestScripts/(run-ship-change|run-ship-signoff|run-ship-change-fail|describe-golden)$' passes. run-ship-change runs the built-in ship-change via run start with fake rex/claude/gh on two leaves (second depends-on first): both built and verified in order, review's two checks ran, PR #7 created once with the run summary as body, run stopped signoff_waiting with no pr merge until rhiza run approve, then land, squash merge pinned to the landed head, main-green and close-out pass; a clone of the pushed PR branch passes rhiza check. TestShipChangeBuiltin asserts the step order with render and render-pr and that ship-change.toml has no shipit; the toml has no ship check --pr-open. Read the diff: render-pr command step sits between review and open-pr with max_visits 3, matching the signoff loop, so request-changes cycles render before each push. Mutation in a scratch clone (review.on.pass straight to open-pr in run-ship-signoff) makes rhiza check on the re-pushed clone fail, so the new check guards the attempt-4 QA finding.

## Links

### Out

- verifies [ship-change-file](../capability/ship-through-rhiza/ship-change-file.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
