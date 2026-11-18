# headless-tabs-index-lock-plan

source · `01M7PSQJ4KXWA4M3043WG83YPC`

## quotes

```json
[]
```

## raw_sha256

9eadd52aa9a2a8829ab3115229683d0aa1cb0e46704e773afd2181e7434baea9

## session

9a5bbf8c-0e64-417e-8304-01a1ed984101

## summary

Autonomous plan run, no owner. Both CI failures were read from the failed first attempts. PR #513 (2026-11-13): leaf b's verifier exited 0 after writing Evidence but the runner judged 'verifier ended with neither Evidence nor findings' (a stale model read); change set 01M7942GRF (verifier-verdict-reads-fresh-model, PR #516) fixed that, and the read path has been rewritten since (#603, #606). PR #593 (2026-11-15 13:51 UTC): leaf a's Evidence commit was refused by update-ref ('cannot lock ref HEAD: is at ...') because the other leaf's builder committed on the shared run worktree's HEAD meanwhile; PR #595 (merged 17:39 UTC the same day) added commitOnly's HEAD-moved retry, covered by TestCommitOnlyHeadMoved. Reproduced at main f1c553613 by running the script 15 and then 20 times in sequence (capped scope, -p 2): one failure in each batch, a third symptom: leaf b's attempt-1 fake builder exited 1 because its git add hit 'Unable to create .../.git/worktrees/repo.built/index.lock: File exists', held by concurrent git work in the shared run worktree (the runner committing leaf a's Evidence). The fake bin/claude builder gives up on the first lock, where a real agent would retry. Planned: fix the fixture so the fake builder retries its add and commit while the index is locked; do not quarantine. Noticed only: the runner's own git calls on the real index (commitOnly's update-index after a successful update-ref) can lose the same lock to a builder's git, filed as a discovered-from draft change set.

## title

Plan of the headless-run-agent-tabs flake (autonomous research)

## Links

### In

- derived-from [headless-tabs-fix-fixture-lock](../decision/headless-tabs-fix-fixture-lock.md)
- derived-from [evidence-commit-waits-index-lock](../capability/evidence-commit-waits-index-lock.md)

## History

- `01M7PSG0YN06CTQCB1Z4KZ77WQ` Fix flaky test TestScripts/headless-run-agent-tabs: failed then passed on a rerun in PR #513 at head 7a96ac202aae8f55b2f0adfcd46d0e0c70574172, PR #593 at head 4c727d1d1f52c4b94be6856422d2b1376284dd54; reproduce it locally, repeatedly, then fix it or quarantine it
