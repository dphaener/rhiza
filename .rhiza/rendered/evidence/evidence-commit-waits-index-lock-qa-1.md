# evidence-commit-waits-index-lock-qa-1

evidence · `01M7PWV0ZYA0JF12F7JSB9Y35F` · actor [claude-code-d5114a0b](../actor/claude-code-d5114a0b.md)

## result

PASS

## statement

QA at 21c89c3b5: fresh go build -trimpath to a scratch path; in a capped scope with RHIZA_/SHIP_/REX_ unset, go test -trimpath -p 2 -count=1 ./internal/api -run 'TestCommitOnly|TestEvidenceCommit' passed (TestCommitOnlyIndexLocked 0.57s, TestCommitOnlyIndexLockedThroughout 2.13s, i.e. failing after the ~2 s bound). By hand with git 2.53 in a scratch linked worktree: a held index.lock makes update-index --add exit 128 with "Unable to create '...index.lock': File exists.", which the retry's match covers; a missing path fails with 'does not exist', which it does not match, so other failures fail at once.

## Links

### Out

- verifies [evidence-commit-waits-index-lock](../capability/evidence-commit-waits-index-lock.md)

## History

- `01M7PSG0YN06CTQCB1Z4KZ77WQ` Fix flaky test TestScripts/headless-run-agent-tabs: failed then passed on a rerun in PR #513 at head 7a96ac202aae8f55b2f0adfcd46d0e0c70574172, PR #593 at head 4c727d1d1f52c4b94be6856422d2b1376284dd54; reproduce it locally, repeatedly, then fix it or quarantine it
