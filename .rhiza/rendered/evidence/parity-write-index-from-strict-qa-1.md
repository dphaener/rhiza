# parity-write-index-from-strict-qa-1

evidence · `01M7GVFM27B77N0BW2BMM0BMNS` · actor [claude-code-4b0e81ca](../actor/claude-code-4b0e81ca.md)

## result

PASS

## statement

QA cycle 1, HEAD build (index path, the CLI default): dry runs with pending entries (set cap-a --dry-run after a pending set and move) give base = the node's last landed op 01M7GV9GCCYS... on disjoint heads and merge the pending update; set --rebase moved a stale update onto the landed version and the change set landed; writes on a landed same-field conflict refuse with conflict exactly as the base build does.

## Links

### Out

- verifies [parity-write-index-from-strict](../capability/rhiza/views/fast-reads/parity-write-index-from-strict.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
