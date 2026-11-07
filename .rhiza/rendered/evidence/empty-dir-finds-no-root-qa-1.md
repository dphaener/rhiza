# empty-dir-finds-no-root-qa-1

evidence · `01M6T6YBQA3Y8624FY1QWNKFFM` · actor [claude-code-960e97dc](../actor/claude-code-960e97dc.md)

## result

PASS

## statement

QA at e6c86be2 with a built binary (go build ./cmd/rhiza): from a subfolder of a scratch dir outside any root, rhiza root show exits 2 'not inside a Rhiza root'; after rhiza root init on the parent, root show from the subfolder (text and --json) names the parent as root; from this checkout it names /home/rhiza/code, project rhiza: the CLI still resolves the root from its absolute cwd. go test -trimpath ./internal/rootcfg/... and ./internal/api/... -run Usage pass (TestFindEmptyDirFindsNoRoot, TestUsageLearnEmptyWdStaysInRepoFolder). In a git-archive copy under a scratch root with rootcfg.go reverted to 2a7a01a1^, the usage tests fail and a fake hit (hit_at 2026-11-04T10:45Z) lands in the scratch root's .rhiza-root/usage.limits.jsonl; with the fix restored they pass and no .rhiza-root is created.

## Links

### Out

- verifies [empty-dir-finds-no-root](../capability/rhiza/runtime/empty-dir-finds-no-root.md)

## History

- `01M6T20EGCJZNFEY6YPJ1PSEDW` Fix tests write fake usage-limit hits into the machine's real root: Ctx.Wd is empty in internal/api tests (testCtx), so limitsDir's rootcfg.Find("") resolves the process cwd, the package dir inside a checkout under the root (~/code), and usage_learn_test's usage-fake hits land in ~/code/.rhiza-root/usage.limits.jsonl; the live root daemon then held every run for a day on a fake weekly hit (and learned fake limits). Root resolution must never escape a test's temp dir: tests must set Wd (or RHIZA_ROOT-style isolation) to their temp dir, and a test must fail if anything is written outside it
