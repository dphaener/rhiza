# verifier-tier-fields

capability · `01M79CZ27B7NJD16PRCG0YJ54R` · parent [build-workflow](../build-workflow.md)

## acceptance

Tests in internal/workflow parse a verifier with both fields, refuse them on a builder, a fixer, a check and an agent step, refuse a blank model and an empty glob, and match `**/testdata/**`, `*_test.go` and `.github/**` against sample paths; cmd/rhiza testscript workflow-describe shows both fields for a verifier that sets them (describe-golden regenerated if the --json shape changes).

## statement

A leaves or checks step's verifier role takes `low_risk_models` (a model ladder by attempt) and `risky_paths` (path globs, `**` matching any number of directories) beside its `models`, which stays the ladder for risky leaves; setting either on a builder, fixer, check or agent step is refused with an error naming the verifier role, a blank model or empty glob is refused, a role setting neither behaves exactly as today, and `rhiza workflow describe` (text and --json) shows both fields.

## Links

### Out

- decided-by [tiered-verifier-rules](../../../../../decision/tiered-verifier-rules.md)

### In

- depends-on [verifier-tier-defaults](verifier-tier-defaults.md)
- verifies [verifier-tier-fields-verified](../../../../../evidence/verifier-tier-fields-verified.md)
- verifies [verifier-tier-fields-qa-1](../../../../../evidence/verifier-tier-fields-qa-1.md)
- verifies [verifier-tier-fields-qa-2](../../../../../evidence/verifier-tier-fields-qa-2.md)
- verifies [pr-525-tiered-verifier](../../../../../evidence/pr-525-tiered-verifier.md)

## History

- `01M76XS745J8S0J74SZFEHQ17Z` Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).
