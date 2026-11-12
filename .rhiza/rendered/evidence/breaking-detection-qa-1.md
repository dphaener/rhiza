# breaking-detection-qa-1

evidence · `01M7836M65BE1HZCS0TAND3XKP` · actor [claude-code-1cc3dbee](../actor/claude-code-1cc3dbee.md)

## result

PASS

## statement

QA at 15f674d4c with a built tools/release binary in a scratch clone: 'release breaking' against the real v0.1.0 tag built both refs in ~12s and printed 'setting: rhiza.daemon.autoUpdate changed type from "off or main" to "off, main, latest or release"' and 'setting: rhiza.daemon.trimGoCache removed' (trimGoCache really is gone from docs/reference/config.md since v0.1.0), '2 breaking difference(s) since v0.1.0', exit 0, no leftover worktree or temp dir. Then tagged HEAD v0.2.0 and committed a renamed show request json field (view -> viewx) and FormatVersion 5 -> 6: it printed 'field: show request.view removed' and 'log-format: version changed from 5 to 6' (the added viewx not reported). go test -trimpath -count=1 ./tools/release/... passes.

## Links

### Out

- verifies [breaking-detection](../capability/rhiza/distribution/breaking-detection.md)

## History

- `01M76B0BSVESMEFD39RTW517R3` Breaking-change detection: compare the API schemas (describe --json), daemon settings and log format at HEAD with the last release tag and report each breaking difference
