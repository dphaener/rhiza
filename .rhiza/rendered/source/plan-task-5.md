# plan-task-5

source · `01M5SWRGRDB40TG5A8G8PGQ2FR`

## quotes

```json
[]
```

## raw_sha256

d9ed4375ac9c0812930c2975abd858731b5cb150866aeefb00b12c47ffcb137f

## summary

The shipit work plan for task 5: a derived SQLite index per worktree, caught up from the log by reading only new entry files, holding the as-is and to-be views, plus a deterministic 50,000-node generator and a read budget test. Decided (D1-D15, agent wording accepted by the owner): a new internal/index package with no command wired and the api's interim model untouched; the pure-Go modernc.org/sqlite driver, so no cgo; the index lives in the user cache folder keyed by a hash of the worktree root; input-incremental catch-up that decodes only new entries, refolds touched change sets and recomputes a view whole when its inputs change, writing rows by diff; as-is and to-be are stored and what-if is computed on demand and benchmarked; freshness is judged by entry file names, and any vanished file forces a full rebuild; the generator is a Go package, deterministic from a seed, with a realistic mix of about 35,000 nodes and 15,000 links over about 500 change sets; a budget test in make check plus benchmarks; a schema version, rebuild on mismatch or corruption, and catch-up under an immediate transaction with WAL and a busy timeout; a minimal per-view query API; tests after implementation, including a property that any split into catch-up steps equals a scratch build; the fresh-read median must stay under 25 ms over 20 runs; no ADR. Revised after the plan: the orchestrator later raised the fresh-read budget to 50 ms, milestone 1's target, after CI on main failed at 27.4 ms (PR #11). Left open: a faster fresh read so CI has headroom under 25 ms again.

## title

Shipit plan, milestone 1 task 5: the SQLite index

## Links

### In

- derived-from [fast-reads](../capability/rhiza/views/fast-reads.md)
- derived-from [t5-d1-index-scope](../decision/t5-d1-index-scope.md)
- derived-from [t5-d2-modernc-sqlite](../decision/t5-d2-modernc-sqlite.md)
- derived-from [t5-d3-index-location](../decision/t5-d3-index-location.md)
- derived-from [t5-d4-input-incremental](../decision/t5-d4-input-incremental.md)
- derived-from [t5-d5-stored-views](../decision/t5-d5-stored-views.md)
- derived-from [t5-d6-freshness-by-names](../decision/t5-d6-freshness-by-names.md)
- derived-from [t5-d7-generator-package](../decision/t5-d7-generator-package.md)
- derived-from [t5-d8-budget-test](../decision/t5-d8-budget-test.md)
- derived-from [t5-d9-index-failures](../decision/t5-d9-index-failures.md)
- derived-from [t5-d10-query-api](../decision/t5-d10-query-api.md)
- derived-from [t5-d11-generator-shape](../decision/t5-d11-generator-shape.md)
- derived-from [t5-d12-tests-after](../decision/t5-d12-tests-after.md)
- derived-from [t5-d13-budget-threshold](../decision/t5-d13-budget-threshold.md)
- derived-from [t5-d14-task-5-guardrails](../decision/t5-d14-task-5-guardrails.md)
- derived-from [t5-d15-no-adr](../decision/t5-d15-no-adr.md)

## History

- `01M5SWRGHKJKT6JZY6ZVMQEGWN` Founding record: the Sources, the root intents and constraints, and the founding decisions
