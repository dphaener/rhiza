# breaking-detection-verified

evidence · `01M78301C7RPJ2DZYG9MPVGD20` · actor [claude-code-f062d0f9](../actor/claude-code-f062d0f9.md)

## result

PASS

## statement

Verifier at 9e4f6f0da: go test -trimpath -count=1 ./tools/release/... passes. TestBreaking's fixture pairs flag a removed op, removed request and response fields, a retyped field, a removed and a retyped setting and a log format bump 5->6 as breaking, and report nothing for an identical surface or for an added op, fields and setting. The diff's parsers match the real inputs: docs/reference/config.md has 46 setting headings and 46 '- Type:' lines, internal/entry/entry.go has 'const FormatVersion = 5', and describe-golden's schemas are inline (no $ref). End to end in a scratch clone tagged v9.9.9 at HEAD~3, 'go run ./tools/release breaking' built both refs, printed '0 breaking difference(s) since v9.9.9', exited 0 and removed its worktree.

## Links

### Out

- verifies [breaking-detection](../capability/rhiza/distribution/breaking-detection.md)

## History

- `01M76B0BSVESMEFD39RTW517R3` Breaking-change detection: compare the API schemas (describe --json), daemon settings and log format at HEAD with the last release tag and report each breaking difference
