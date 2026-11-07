# tests-leave-root-state-alone-qa-1

evidence · `01M6T6YBS2HXNPQB1WEKHDYBMX` · actor [claude-code-960e97dc](../actor/claude-code-960e97dc.md)

## result

PASS

## statement

QA at e6c86be2: go test -trimpath ./internal/rootguard/... passes. From this checkout under the real root ~/code, ./internal/api/... -run Usage, the full ./internal/api/... and ./cmd/rhiza -run TestScripts/usage ran with a before/after snapshot of ~/code/.rhiza-root (file list, sizes, mtimes, limits file hash): identical, nothing written. Negative check: in a copy under a scratch root with the rootcfg fix reverted, rootguard.Main in api's TestMain failed the run with '.rhiza-root/usage.limits.jsonl: line 1 appended with hit_at 2026-11-04T10:45:00Z, outside the test run'. The full api run had one unrelated failure, TestWebRootHandlerServesEveryProject's TempDir cleanup ('directory not empty'), which passes 3/3 alone and is in code this change does not touch.

## Links

### Out

- verifies [tests-leave-root-state-alone](../capability/rhiza/integrity/self-check/tests-leave-root-state-alone.md)

## History

- `01M6T20EGCJZNFEY6YPJ1PSEDW` Fix tests write fake usage-limit hits into the machine's real root: Ctx.Wd is empty in internal/api tests (testCtx), so limitsDir's rootcfg.Find("") resolves the process cwd, the package dir inside a checkout under the root (~/code), and usage_learn_test's usage-fake hits land in ~/code/.rhiza-root/usage.limits.jsonl; the live root daemon then held every run for a day on a fake weekly hit (and learned fake limits). Root resolution must never escape a test's temp dir: tests must set Wd (or RHIZA_ROOT-style isolation) to their temp dir, and a test must fail if anything is written outside it
