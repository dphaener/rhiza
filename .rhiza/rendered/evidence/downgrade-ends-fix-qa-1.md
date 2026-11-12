# downgrade-ends-fix-qa-1

evidence · `01M78XSAAS8RFPQY80J140RWGX` · actor [claude-code-9af4801e](../actor/claude-code-9af4801e.md)

## result

PASS

## statement

QA at e313e4634 with the built binary, a real review run with checks alpha and beta (fake agents). Alpha's only major was downgraded to minor: alpha passed in cycle 1 ('passed: critic downgraded f1 (major to minor)'), and the session log shows no alpha re-run and no fixer for it. Beta's blocker was downgraded to major and still started the fixer and verifier for leaf thing in cycle 2, then beta re-ran and passed. rhiza run status shows 'f1 major ... downgraded from major to minor', and --json keeps severity major with downgraded_to minor. rhiza run summary lists both under 'Downgraded by the critic' as 'downgraded from <orig> to <lower>', not under Confirmed.

## Links

### Out

- verifies [downgrade-ends-fix](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-critic/downgrade-ends-fix.md)

## History

- `01M76XRYNK4PGG53XNC1ZWFTKD` The review critic can confirm, downgrade or drop a blocker or major finding; a downgraded finding keeps the lower severity, so a real but minor defect no longer starts a fix cycle; review and plan review share one severity contract with a per-check mapping. Evidence: critics confirmed 242 of 246 blocker/major findings; 57% of runs needed a second review cycle (prompt-gap rec 5).
