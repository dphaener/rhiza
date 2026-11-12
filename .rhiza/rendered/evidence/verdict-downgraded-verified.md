# verdict-downgraded-verified

evidence · `01M78WYPPHD6N9HXTPYGK22735` · actor [claude-code-18916c15](../actor/claude-code-18916c15.md)

## result

PASS

## statement

At 5c8b9d932, verifier ran go test -trimpath ./internal/api -run RunVerdict -v: TestRunVerdictDowngrade PASS. It journals major->minor and blocker->major with severity and original_severity, and refuses as CodeUsage a downgrade of a major to major, a major or blocker to blocker, one with no severity, one of a minor finding, and a severity on a confirmed verdict. go test -trimpath ./internal/run: ok (TestStandingDowngrade covers the journaled verdict and that a finding downgraded below major no longer stands). go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden: ok, and re-running it with -update left no diff, so the golden is current. Diff read: SeverityRank orders blocker<major<minor<nit, and the validation requires rank(new) > rank(original). Note: run.Verdict now decodes original_severity from input, so a confirmed or dropped verdict could carry a caller-supplied one (a downgraded verdict's is overwritten). This is outside the acceptance.

## Links

### Out

- verifies [verdict-downgraded](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-critic/verdict-downgraded.md)

## History

- `01M76XRYNK4PGG53XNC1ZWFTKD` The review critic can confirm, downgrade or drop a blocker or major finding; a downgraded finding keeps the lower severity, so a real but minor defect no longer starts a fix cycle; review and plan review share one severity contract with a per-check mapping. Evidence: critics confirmed 242 of 246 blocker/major findings; 57% of runs needed a second review cycle (prompt-gap rec 5).
