# review-guide-fakes-only

capability · `01M6505MN5E02V3Y30RMC16FN5` · parent [review-guide](../review-guide.md)

## acceptance

go test -trimpath ./internal/guide/... passes with a test that finds the fakes-only rule, naming Rex, agents and gh, in the expanded check-compliance, check-quality and check-model topics, finds in check-qa that qa is hands-on with the built binary and that fakes alone do not count, finds no scope check named, and that each topic stays within TopicCap

## statement

The check-compliance, check-quality and check-model topics tell their checks to test runs with fakes only, never real Rex, agents or gh; check-qa runs the built binary hands-on, since fakes alone do not count for it

## Links

### Out

- decided-by [fakes-only-in-build-and-review](../../../../../../decision/fakes-only-in-build-and-review.md)

### In

- verifies [review-guide-fakes-only-verified](../../../../../../evidence/review-guide-fakes-only-verified.md)
- verifies [review-guide-fakes-only-qa-1](../../../../../../evidence/review-guide-fakes-only-qa-1.md)
- verifies [pr-104-agents-never-start-runs](../../../../../../evidence/pr-104-agents-never-start-runs.md)
- verifies [pr-537-review-role-topics](../../../../../../evidence/pr-537-review-role-topics.md)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.
- `01M76XRETZFTD2TQ5GS3WYAERA` Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).
