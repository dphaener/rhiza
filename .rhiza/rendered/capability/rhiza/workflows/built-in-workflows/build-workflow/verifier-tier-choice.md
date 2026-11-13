# verifier-tier-choice

capability · `01M79CZ29TQH0XYRCP38C2X9DK` · parent [build-workflow](../build-workflow.md)

## acceptance

Tests in internal/run start verifiers from a role with low_risk_models = ["sonnet"], models = ["opus"] and risky_paths = [".github/**"] and find the journalled leaf-session model sonnet for a first-attempt commit touching only internal/x/x.go and for review cycle 2's first fix touching only that file, and opus for a commit touching .github/workflows/ci.yml, for attempt 2, for a fresh round, for a fix retry within review cycle 2, and when the start HEAD cannot be diffed; a review fix in cycle 3 is decided risky whatever the diff; a role without low_risk_models starts on models as before.

## statement

When the conductor launches a verifier whose role sets `low_risk_models`, it starts it on that ladder only when the leaf is low-risk, and on `models` when the leaf is risky: a file changed between the builder's or fixer's journalled start HEAD and the commit under verification matches one of the role's `risky_paths`, the attempt within the round is past 1, the leaf is in a fresh round, on a checks step the review cycle is past 2 or the fix try is past the first of its cycle (fixes start in cycle 2), or the diff cannot be read; the choice is deterministic, with no agent call, and the leaf-session entry journals the model chosen.

## Links

### Out

- decided-by [tiered-verifier-rules](../../../../../decision/tiered-verifier-rules.md)

### In

- verifies [verifier-tier-choice-verified](../../../../../evidence/verifier-tier-choice-verified.md)
- verifies [verifier-tier-choice-qa-1](../../../../../evidence/verifier-tier-choice-qa-1.md)
- verifies [verifier-tier-choice-verified-2](../../../../../evidence/verifier-tier-choice-verified-2.md)
- verifies [verifier-tier-choice-qa-2](../../../../../evidence/verifier-tier-choice-qa-2.md)
- verifies [pr-525-tiered-verifier](../../../../../evidence/pr-525-tiered-verifier.md)

## History

- `01M76XS745J8S0J74SZFEHQ17Z` Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).
