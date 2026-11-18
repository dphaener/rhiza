# api-tests-tmp-leak-guard-qa-1

evidence · `01M7Q2DZGY9KDZC540BA875J04` · actor [claude-code-91a691bb](../actor/claude-code-91a691bb.md)

## result

PASS

## statement

QA: acceptance go test passes, scratch TMPDIR empty after (no rhiza-api-tmp- left), 0 new /tmp/rhiza-projection-*; an injected rhiza-projection-qa in testTmp makes the binary print its name and exit 1, folder still removed; a 2h-old rhiza-api-tmp- folder is swept, a fresh one kept

## Links

### Out

- verifies [api-tests-tmp-leak-guard](../capability/rhiza/views/fast-reads/api-tests-tmp-leak-guard.md)

## History

- `01M7PZQ416EJR49M7DRVEBJ8WH` Fix Throwaway projections leak into /tmp and fill RAM: since #606 (one read path, throwaway-index-fallback) test runs leave rhiza-projection-<n>/ folders (index.db, -shm, -wal, 0.3-3 MB each) under the real /tmp and never remove them. On 2026-11-17, 8,900 of them (7.1 GB of a RAM-backed /tmp tmpfs) piled up at roughly 1,000 per hour during builds. With the VM ballooned to about 5.7 GB of RAM and swap full, earlyoom SIGTERMed 'rhiza ship open-pr' twice (VmRSS 2.2 GB each, 06:25Z and 06:26Z), and plan-bug run 01M7M27HT2D61908ZS667F197T failed at plan/open-pr (exit 143). None appear while no test runs, so the leak is in tests (likely tests that open a projection with no IndexPath, outside testscript's per-script TMPDIR), possibly also in production fallbacks. Expected: every throwaway projection's folder is removed when it closes, tests put theirs under t.TempDir(), and the test leak guard checks that no folder is left. Also look at why 'rhiza ship open-pr' needs 2.2 GB RSS.
