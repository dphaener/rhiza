# breaking-detection-verified-2

evidence · `01M783HYY0W9FTVC2YVHPXB3CG` · actor [claude-code-aae3425f](../actor/claude-code-aae3425f.md)

## result

PASS

## statement

Verifier, attempt 2 at ca38b73c6: go test -trimpath -count=1 -v ./tools/release/... passes. TestBreaking: removed op, removed request and response field, retyped field, removed and retyped setting and a log format bump are each reported as breaking; added op, field and setting give none. TestBreakingOneOfAndMaps: a field removed or retyped inside a oneOf variant, a dropped variant and a retyped map value are breaking; reordered or added variants are not. Read the diff: input documents, oneOf variants (matched by const properties) and additionalProperties are now decoded and compared, which fixes quality f1. Ran Breaking on this build's real describe --json: against itself it reports nothing; with handle removed from apply's set write variant it reports 'field: apply input.writes[](op=set).handle removed'. The minor findings left open (an enum superset is reported as a retyped setting, a newly required field and a removed enum value are not detected, the v* tag pattern, the breaking subcommand has no test) are outside the leaf's acceptance.

## Links

### Out

- verifies [breaking-detection](../capability/rhiza/distribution/breaking-detection.md)

## History

- `01M76B0BSVESMEFD39RTW517R3` Breaking-change detection: compare the API schemas (describe --json), daemon settings and log format at HEAD with the last release tag and report each breaking difference
