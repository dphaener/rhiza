# breaking-detection

capability · `01M76B7DRSHQHDGWY46MBH319P` · parent [distribution](../distribution.md)

## acceptance

go test -trimpath ./tools/release/... over fixture pairs: removed op, removed field, retyped field, removed setting and a log format bump are breaking; an added op, field or setting is not

## statement

tools/release lists each breaking difference between HEAD and the last v* tag: a removed operation or request/response field, or a field whose type changed, in describe --json; a removed or retyped daemon setting; a changed log format version

## Links

### Out

- decided-by [m33-breaking-is-major](../../../decision/m33-breaking-is-major.md)

### In

- verifies [breaking-detection-verified](../../../evidence/breaking-detection-verified.md)
- verifies [breaking-detection-qa-1](../../../evidence/breaking-detection-qa-1.md)
- verifies [breaking-detection-verified-2](../../../evidence/breaking-detection-verified-2.md)
- verifies [pr-493-breaking-change-detection](../../../evidence/pr-493-breaking-change-detection.md)

## History

- `01M76B0BSVESMEFD39RTW517R3` Breaking-change detection: compare the API schemas (describe --json), daemon settings and log format at HEAD with the last release tag and report each breaking difference
