# headless-tabs-builder-waits-index-lock-verified

evidence · `01M7PW907QAFR33NQZZ21SVJZ6` · actor [claude-code-4986ec52](../actor/claude-code-4986ec52.md)

## result

PASS

## statement

At 21c89c3b5, in capped scopes with RHIZA_/SHIP_/REX_ unset: (1) go test -trimpath -p 2 -count=30 ./cmd/rhiza -run 'TestScripts/^headless-run-agent-tabs$' exit 0 (52s). (2) Scratch, uncommitted edit to the fake rex creating $WORK/repo/.git/worktrees/repo.built/index.lock before each builder pane and removing it 0.3s later: with the old builder (single git add/commit) the same command at -count=1 exits 1, all six builder attempts 'fail (exit 1)'; with the new retrying builder it exits 0; scratch edit reverted, cmd/ clean. (3) The diff touches only the txtar's header comment and the bin/claude builder body; no assertion lines changed, no Go code, no other scripts. rhiza ship generate and rhiza render exit 0.

## Links

### Out

- verifies [headless-tabs-builder-waits-index-lock](../capability/headless-tabs-builder-waits-index-lock.md)

## History

- `01M7PSG0YN06CTQCB1Z4KZ77WQ` Fix flaky test TestScripts/headless-run-agent-tabs: failed then passed on a rerun in PR #513 at head 7a96ac202aae8f55b2f0adfcd46d0e0c70574172, PR #593 at head 4c727d1d1f52c4b94be6856422d2b1376284dd54; reproduce it locally, repeatedly, then fix it or quarantine it
