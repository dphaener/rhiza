# ship-one-ci-rebase-in-steps-verified

evidence · `01M7NZDMM9ZZ47X914FS1QTZ1B` · actor [claude-code-ddbf036f](../actor/claude-code-ddbf036f.md)

## result

PASS

## statement

In a detached checkout at 337430c11, go test -p 2 -trimpath ./cmd/rhiza -run 'TestScripts/ship-change-one-ci$' -count=1 (capped scope, RHIZA_/SHIP_/REX_ unset) exit 0. The script asserts: happy path wait-ci then pr-green then merge, with no rebase, no wait-ci (visit 2) and no handback. Main moved cleanly: merge (visit 1) fail then rebase, rebase (visit 1) pass then wait-ci, wait-ci/pr-green/merge (visit 2) pass, conduct visit 1 only, no signoff (visit 2), no handback in status or claude.log, 2 pushes, 1 pr merge, landed. Failure case (CONFLICT, main adds first.txt): rebase (visit 1) fail then conduct, conduct (visit 2) pass then gate, the fake conductor resolves a MERGE_HEAD merge (resolved logged once), then gate, signoff, push (visit 2) and merge (visit 2) pass, 1 pr merge, landed, origin main first.txt=first. rhiza ship generate exit 0 (0 changed); rhiza render exit 0, apart from one actor-line ordering drift in index.md that was already present before this commit. The diff touches only the txtar plus the builder's log and rendered files: no ship-change.toml, no internal/api, no processScripts change (the script was already listed for fake rex/claude and sh -c), the fake claude now runs only git, no network or real gh, no TODOs.

## Links

### Out

- verifies [ship-one-ci-rebase-in-steps](../capability/rhiza/orchestration/shipping/ship-one-ci-rebase-in-steps.md)

## History

- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict
