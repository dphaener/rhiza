# check-actor-records-committed

capability · `01M6504QPJR31GMTNTV530MEHJ`

## acceptance

TestCommitCheckWritesNone is changed: a check or critic session whose actor wrote no model entries, but whose actor's record entries are uncommitted, makes a commit named record <check> <role> writes with Change: and Actor: trailers and no Leaf: trailer, and afterwards the run worktree has no uncommitted record entries for that actor; a session whose actor has nothing uncommitted still makes no commit

## statement

A review check or critic session that writes no model entries still gets its actor's record entries committed, so no uncommitted record change set is left in the run worktree for a later fixer to sweep up

## Links

### Out

- decided-by [commit-actor-record-entries](../decision/commit-actor-record-entries.md)

### In

- verifies [check-actor-records-committed-verified](../evidence/check-actor-records-committed-verified.md)
- verifies [check-actor-records-committed-qa-1](../evidence/check-actor-records-committed-qa-1.md)
- verifies [pr-99-commit-idle-check-actors](../evidence/pr-99-commit-idle-check-actors.md)

## History

- `01M602TH4C724WZB9MKBN5Z1PP` A review check or critic session that writes no model entries still gets its actor's record entries committed, so no uncommitted record change set is left in the run worktree for a later fixer to sweep up
