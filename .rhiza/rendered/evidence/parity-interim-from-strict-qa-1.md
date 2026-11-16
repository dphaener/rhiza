# parity-interim-from-strict-qa-1

evidence · `01M7GVFKVRH96QV0JYR496RT6E` · actor [claude-code-4b0e81ca](../actor/claude-code-4b0e81ca.md)

## result

FAIL

## statement

QA cycle 1, HEAD build in a scratch repo with the default index: disjoint heads (two landed updates on statement and acceptance based on one create) show merged fields, version the last applied op, and new set/mv/retire name it as base and land; a draft's child under a parent retired by a later landed set is absent (show/set: not_found, check warns retired parent). ready, overlaps, milestones, ls and daemon status match the base build. Broken: rhiza inbox --json and health drop the accept item of a proposed change set whose leaf has a landed parent (1 of 3 listed; base lists 3); see finding f1.

## Links

### Out

- verifies [parity-interim-from-strict](../capability/rhiza/views/fast-reads/parity-interim-from-strict.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
