# ship-one-ci-rebase-in-steps-qa-2

evidence · `01M7P1AFSHMVDTMYH899KNS2T9` · actor [claude-code-1bba01de](../actor/claude-code-1bba01de.md)

## result

PASS

## statement

QA: go test -p 2 -trimpath -count=1 -v ./cmd/rhiza -run 'TestScripts/(ship-change-one-ci|guide|workflow-new)$' (capped scope, env cleaned) PASS, ship-change-one-ci 4.99s: happy path one CI wait/merge, main-moved case merge(visit 1) fail then rebase pass then wait-ci/pr-green/merge visit 2 with conduct visited once and no handback; conflict case rebase fail then conduct (visit 2), gate, signoff, push, merge visit 2.

## Links

### Out

- verifies [ship-one-ci-rebase-in-steps](../capability/rhiza/orchestration/shipping/ship-one-ci-rebase-in-steps.md)

## History

- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict
