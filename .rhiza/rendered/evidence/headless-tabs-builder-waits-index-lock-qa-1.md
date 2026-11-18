# headless-tabs-builder-waits-index-lock-qa-1

evidence · `01M7PWV11V9FYDYJSZMASTCZ1Y` · actor [claude-code-d5114a0b](../actor/claude-code-d5114a0b.md)

## result

PASS

## statement

QA at 21c89c3b5, capped scope, RHIZA_/SHIP_/REX_ unset: go test -trimpath -p 2 -count=10 ./cmd/rhiza -run 'TestScripts/^headless-run-agent-tabs$' passed (18.5s). With a temporary uncommitted edit making the fake rex create the pane worktree's index.lock (git rev-parse --git-path) before each pane and remove it 0.3 s later: -count=3 passed with the new fake builder; with the old fake builder put back it failed with 'attempt N builder ... fail (exit 1)'. Edit reverted; the script is unchanged.

## Links

### Out

- verifies [headless-tabs-builder-waits-index-lock](../capability/headless-tabs-builder-waits-index-lock.md)

## History

- `01M7PSG0YN06CTQCB1Z4KZ77WQ` Fix flaky test TestScripts/headless-run-agent-tabs: failed then passed on a rerun in PR #513 at head 7a96ac202aae8f55b2f0adfcd46d0e0c70574172, PR #593 at head 4c727d1d1f52c4b94be6856422d2b1376284dd54; reproduce it locally, repeatedly, then fix it or quarantine it
