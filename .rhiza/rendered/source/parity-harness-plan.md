# parity-harness-plan

source · `01M7GGD9KF03S0ZM6121WB1WC1`

## quotes

```json
[]
```

## raw_sha256

55ad21bafe1e3c7221ffb4560ab8443c207adee90e18a938caebfe5c6cf41855

## session

340e46bd-dcba-4ca7-8898-01a1e104370b

## summary

Autonomous plan run of the Parity harness change set (milestone 39, 2026-11-15). No owner interview: the orchestrator, as arbiter, said design choices are the planner's, to be recorded as assumed Decisions with a reasoning field, and that the owner's 'Ok' (milestone chat, question 2) makes the strict engine in internal/model the spec. The change set already held three leaves from the milestone plan. The research pass found what they left open: the two api models and the strict engine define a node's version differently. The interim model (internal/api/model.go) and write_indexed.go's versionOf take the largest head of the base graph and set Conflict when there are two heads; the strict engine's version is the last operation applied, and a conflict is a Problem only when a stale operation overlaps an aspect (disjoint edits merge silently). Writes refuse a node or link with Conflict set (write_node.go writableNode and buildRetire, write_link.go), and set.txtar and write_test.go LandedConflict pin that for two heads writing the same field, which the strict rule also reports. Pending has no strict counterpart. The index already computes the strict what-if (index.WhatIfOverlay over model.WhatIfDelta), while writeIndex.buildModel re-applies in-flight ops by the interim rules. The interim model is deleted by the milestone's projection change set. newModelSets is also called per candidate in daemon_fold.go, proposals.go and inbox.go. The planner settled the version and conflict rule the harness compares, how disagreements are fixed (derive from the strict engine instead of re-implementing it), that the harness lands green with a named known-cases table, and that the generator of internal/index/property_test.go moves to a test-only package. It refined the three leaves to say so.

## title

Parity harness plan, autonomous

## Links

### In

- derived-from [parity-version-rule](../decision/parity-version-rule.md)
- derived-from [parity-fix-derives-from-strict](../decision/parity-fix-derives-from-strict.md)
- derived-from [parity-known-cases-table](../decision/parity-known-cases-table.md)
- derived-from [parity-generator-shared](../decision/parity-generator-shared.md)

## History

- `01M7GFV2ZMHRJR78H0YTMKWXB7` A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.
