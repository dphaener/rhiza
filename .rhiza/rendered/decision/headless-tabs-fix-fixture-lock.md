# headless-tabs-fix-fixture-lock

decision · `01M7PSQZKR8DG3J4VXDYVTNPJK`

## statement

Fix the headless-run-agent-tabs flake in its fixture, not by quarantine: the fake bin/claude builder retries its git add and git commit while the shared run worktree's index is locked by concurrent git work, as a real agent would; the two CI symptoms of PRs #513 and #593 are already fixed (verifier-verdict-reads-fresh-model, and commitOnly's HEAD-moved retry from PR #595), and the index.lock race is the one that still reproduces at main

## wording

Fix the fixture's index.lock race; do not quarantine. Planner's ruling in an autonomous plan run, from the reproduction

## Links

### Out

- derived-from [headless-tabs-index-lock-plan](../source/headless-tabs-index-lock-plan.md)

### In

- decided-by [headless-tabs-builder-waits-index-lock](../capability/headless-tabs-builder-waits-index-lock.md)

## History

- `01M7PSG0YN06CTQCB1Z4KZ77WQ` Fix flaky test TestScripts/headless-run-agent-tabs: failed then passed on a rerun in PR #513 at head 7a96ac202aae8f55b2f0adfcd46d0e0c70574172, PR #593 at head 4c727d1d1f52c4b94be6856422d2b1376284dd54; reproduce it locally, repeatedly, then fix it or quarantine it
