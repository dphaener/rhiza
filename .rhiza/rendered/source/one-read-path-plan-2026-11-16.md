# one-read-path-plan-2026-11-16

source · `01M7HCBXV5EBNYVYR47NG63JB1`

## quotes

```json
[]
```

## raw_sha256

194309b8937924864952399ec88777136b0a1da71f764b4eadbf7d54ec5033f1

## session

6072e9ae-0ea5-46d3-a745-01a1e2c2427e

## summary

Autonomous plan run (no owner). The arbiter's brief: plan the change set after #599, #600 and #603 landed; inbox (531 ms) and check (7.1 s) are still over projection-targets' 200 ms and 5 s, so add a leaf per command that profiles and fixes it if deleting the interim path will not; stay out of cmd/rhiza's script harness, which change set 01M7GFTS0N rebuilds next. The research pass found that #603 already gave writes the strict what-if view (write_indexed.go uses WhatIfOverlayOver; the parity table is empty) and gave the daemon's projects, question answers at a gate and inbox's gate read their root's own index, each falling back to an empty IndexPath. What remains is every branch that reads the log when there is no projection (workspace.model, entriesOfSet and changeSets; read.go, change.go, ship.go, slice.go, run_model_cache.go, run_milestone.go; status, inbox, proposals, the daemon fold and render) and web_actions.go's cleared IndexPath. One profiled probe on the main checkout at c84997915: rhiza inbox 0.63 s, about half of it run.Store.List decoding all 407 run journals (45 MB); rhiza check 7.35 s, of which model.WhatIfFrom's per-change-set seed copy (copyNode/copyFields) is about 2 s CPU, check.badLinks/inFlightLinks 1.4 s and the entry decode 1.4 s. Neither cost is in the interim model, so the plan adds an inbox leaf and a check leaf. The plan revises the four leaves written with the milestone to this state, adds a leaf for the whole-log folds and moves the throwaway index into internal/api's projection open. No question was escalated; each choice is recorded as an arbiter-worded Decision.

## title

One read path: autonomous plan, 2026-11-16

## Links

### In

- derived-from [throwaway-index-in-api](../decision/throwaway-index-in-api.md)
- derived-from [render-keeps-strict-log-fold](../decision/render-keeps-strict-log-fold.md)
- derived-from [log-read-count-hook](../decision/log-read-count-hook.md)

## History

- `01M7GFTKCAK0GY3PPDRTPVYF56` The interim model is deleted: writes take their model and base versions from the projection under the strict rules, api/model.go's newModel, newModelSets and their apply rules, Ctx.readLog and every IndexPath branch go, and a worktree with no usable index cache reads through a throwaway index, so there is one read path and one set of model rules. Evidence: about a dozen ops branch on IndexPath (write.go:189, read.go:180/197/222, inbox.go:602, ship.go:359, run_model_cache.go:88, repair.go:161) and three clear it (daemon_projects.go:46, question.go:182, web_actions.go:108).
