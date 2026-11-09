# pr-405-fix-flaky-test-on-main

evidence · `01M6Z4WMM68G2S3947KTH25YSH` · actor [dphaener](../actor/dphaener.md)

## head

ec90fb941e654bbe8d37a391eff305912954f1ef

## pr

https://github.com/dphaener/rhiza/pull/405

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #405 (https://github.com/dphaener/rhiza/pull/405) passed CI at ec90fb9

## status

verified

## Links

### Out

- verifies [web-tests-drain-health-refresh](../capability/rhiza/integrity/self-check/web-tests-drain-health-refresh.md)

## History

- `01M6YKX97C202RDFNV5DGVS033` Fix Flaky test on main: TestWebRootHandlerServesEveryProject (internal/api) fails about 1 run in 10 with 'TempDir RemoveAll cleanup: directory not empty'. Reproduced locally on an unmodified git archive of origin/main with go test -trimpath -count=10 -run TestWebRootHandlerServesEveryProject ./internal/api/ (found while verifying the #398 merge, 2026-11-08). Likely a goroutine or background process (web handler, daemon state write, or git) still writing into t.TempDir() after the test returns. Fix the root cause so -count=50 passes; don't add retries.
