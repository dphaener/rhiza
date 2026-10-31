# ci-rerun-once-per-head

decision · `01M6A0MFAVA728T8Y6TNPJVJZQ`

## statement

When a PR or main check fails, ship check reruns that GitHub Actions run's failed jobs once per head (run id from the check's detailsUrl, or databaseId for --main), records the rerun in the ship record and keeps waiting; a second failure on the same head fails as today; checks with no Actions run id are not rerun

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m10-orchestrator-picks](../source/m10-orchestrator-picks.md)

## History

- `01M6A0JENBMBYPYRN5195AGPKH` Plan milestone 10: runs heal themselves, picked by the orchestrator under the owner's standing instruction
