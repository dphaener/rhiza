# unique-finding-ids-verified

evidence · `01M64TJJPEP2CH57GJ2VQJY99R` · actor [claude-code-6de6bf58](../actor/claude-code-6de6bf58.md)

## result

PASS

## statement

At dca2fdd, cmd/rhiza/testdata/script/findings-concurrent.txtar runs four concurrent --leaf findings calls (verifier burst) and four concurrent --check alpha calls, and findings-concurrent-rounds.txtar (added for attempt 1's finding) runs four concurrent plain `rhiza run findings` calls on one run (the plan review round path). They assert the journalled IDs are unique within each scope: leaf f1..f4, check rounds 1..4 each with f1,f2, plain rounds 1..4 each with f1,f2 exactly once, and run status shows 4 rounds of 2. Ran: go test ./cmd/rhiza -run TestScripts/findings-concurrent -count=5 passed. In a scratch worktree with internal/ reverted to the pre-fix 974f0ca, both scripts fail: duplicate leaf f1/f2, and all four plain calls report round 1. go test ./... passed. make check's fmt-check, vet and tests passed, except TestReadBudget, a timing test that fails only under machine load and is outside this diff; self-check reported only the stale render of the builder's actor record. Read the diff: lockJournal in internal/run/lock_unix.go wraps round/ID numbering in runRunFindings (internal/api/run_review.go).

## Links

### Out

- verifies [unique-finding-ids](../capability/unique-finding-ids.md)

## History

- `01M5ZNV40RB8VYBB0A3CPXCNXZ` rhiza run findings numbers findings under the run's lock, so concurrent calls in one session never journal two findings with the same ID (leaf findings and plan review findings alike).
