# build-guide-fakes-only

capability · `01M6505JJY04DWHQZ3EQRFCHAM` · parent [build-guide](../build-guide.md)

## acceptance

go test -trimpath ./internal/guide -run 'TestVerifierText|TestSharedTextsWithinCap' finds the fakes-only rule in the verifier topic, naming Rex, agents and gh, and the topic stays within TopicCap

## statement

rhiza guide verifier (the rhiza-verifier skill) tells verifiers to test runs with fakes only, never real Rex, agents or gh

## Links

### Out

- decided-by [fakes-only-in-build-and-review](../../../../../../decision/fakes-only-in-build-and-review.md)

### In

- verifies [build-guide-fakes-only-verified](../../../../../../evidence/build-guide-fakes-only-verified.md)
- verifies [build-guide-fakes-only-qa-1](../../../../../../evidence/build-guide-fakes-only-qa-1.md)
- verifies [pr-104-agents-never-start-runs](../../../../../../evidence/pr-104-agents-never-start-runs.md)
- verifies [pr-521-build-role-topics](../../../../../../evidence/pr-521-build-role-topics.md)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.
- `01M76XR5Y1N5E699AT2WMJMQJC` Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).
