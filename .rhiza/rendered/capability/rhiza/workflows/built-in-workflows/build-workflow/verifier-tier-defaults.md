# verifier-tier-defaults

capability · `01M79CZ2BQEZSPC3KEF7PZH5G0` · parent [build-workflow](../build-workflow.md)

## acceptance

A test in internal/workflow loads the built-in build and review workflows and finds those ladders and globs on both verifier roles; cmd/rhiza testscript workflow-describe shows them for `rhiza workflow describe build`; `make self-check` passes.

## statement

The built-in build workflow's verifier and review workflow's fixed-leaf verifier keep models = ["opus"] and ship low_risk_models = ["sonnet"] and risky_paths covering test infrastructure (*_test.go, **/*_test.go, **/testdata/**, **/main_test.go), hooks (.githooks/**, **/hooks/**), CI (.github/**), workflow files (.rhiza/workflows/**, internal/workflow/builtin/**), migrations (**/migrations/**, **/migrate/**) and retry code (**/*retry*); a project changes them by overriding the workflow, with no new rhiza.* setting.

## Links

### Out

- depends-on [verifier-tier-fields](verifier-tier-fields.md)
- decided-by [tiered-verifier-rules](../../../../../decision/tiered-verifier-rules.md)

### In

- verifies [verifier-tier-defaults-verified](../../../../../evidence/verifier-tier-defaults-verified.md)
- verifies [verifier-tier-defaults-qa-1](../../../../../evidence/verifier-tier-defaults-qa-1.md)
- verifies [verifier-tier-defaults-verified-2](../../../../../evidence/verifier-tier-defaults-verified-2.md)
- verifies [verifier-tier-defaults-qa-2](../../../../../evidence/verifier-tier-defaults-qa-2.md)
- verifies [verifier-tier-defaults-verified-3](../../../../../evidence/verifier-tier-defaults-verified-3.md)
- verifies [pr-525-tiered-verifier](../../../../../evidence/pr-525-tiered-verifier.md)

## History

- `01M76XS745J8S0J74SZFEHQ17Z` Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).
