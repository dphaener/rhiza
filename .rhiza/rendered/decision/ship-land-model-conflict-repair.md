# ship-land-model-conflict-repair

decision · `01M77NFJV6S5JN933PSDTDX60A`

## reasoning

The question records run 01M77DW ending with green PR #480 unmerged because land was refused on change set 01M763Z4's stale-base update to ship-check-ci-test-findings, and land-triage (rhiza ship gate --failed) only passes for a recorded pre-push failure, so any other land failure ends the run. The model already supports the repair: buildUpdate in internal/api/write_node.go has --rebase, which moves a pending update onto the landed version only when the call restates every field it writes, so landed work is not silently replaced. ship-change already pairs every other refusal class with a triage step and a bounded fixer (rebase-triage -> conflict-fix for code conflicts, gate-fix for gates, ci-fix for CI), so a model conflict should follow the same shape rather than strand a merge-ready PR for the watchdog. Making the gate fail on the shipping change set's own conflicts closes the gap the question names: rhiza check only warns for in-flight conflicts, so every check passed and the problem surfaced only at land. Restricting the hard failure to the shipping change set keeps other in-flight change sets' conflicts from blocking unrelated runs.

## statement

ship-change treats a land refused on a model conflict (the shipping change set's own operation resting on a stale base) as a fixable step, not an end. ship land reports the refusal distinctly (the conflicting nodes and fields, both versions) and land-triage routes it to a model-conflict fixer agent; the fixer restates each conflicting update on the landed version with `rhiza set <ref> --rebase`, keeping the landed change's work and adding this change set's intent, commits the rewritten .rhiza/log entry with the Change/Actor/Review trailers, and the run goes back through the gate and land (bounded by max_visits, so a repeat conflict still fails the run with the conflict named). Earlier, the pre-land gate (and rhiza ship check) fails, rather than only warns, on a conflict in the shipping change set's own operations, so the stale base is repaired before sign-off where possible. Conflicts of other in-flight change sets stay warnings.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-6eec8236-1](../question/question-6eec8236-1.md)
- derived-from [plan-ship-land-model-conflict-repair](../source/plan-ship-land-model-conflict-repair.md)

### In

- relates-to [ship-land-model-conflict-two-designs](../question/ship-land-model-conflict-two-designs.md) — The question asks which of two contradictory designs for a land refused on a model conflict ships; this Decision is one of them.
- supersedes [merge-keeps-main-model-conflict-design](merge-keeps-main-model-conflict-design.md) — The merge keeps main's landed design instead of this change set's.

## History

- `01M77N83F027XBGGZPK9SV6YN2` Fix watchdog unmerged PR at step land-triage: run 01M77DWNRDSFY0PEP0S1P58P5W; log: PR green at 52ace47 and unmerged after step land-triage failed: step "land-triage" ended with outcome "fail" | exit 1 | No head the local gate failed on after a rebase is on the branch.
