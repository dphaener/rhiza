# afk-prompts-shipped-default

decision · `01M76XMNKREKEH8K8E1ETGVR0H`

## quote

Those should ship as the default build tools, and we also want to use them as our own.

## statement

afk's shipit/forge prompts become Rhiza's shipped default prompts for every run role, and Rhiza's own repo runs those defaults

## Links

### Out

- derived-from [efficient-build-cycle-interview](../source/efficient-build-cycle-interview.md)

### In

- decided-by [m35-afk-prompts](../intent/milestone-35/m35-afk-prompts.md)
- decided-by [builder-topic](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/builder-topic.md)
- decided-by [verifier-topic](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/verifier-topic.md)
- decided-by [build-prompts-one-line](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/build-prompts-one-line.md)
- decided-by [check-compliance-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/check-compliance-topic.md)
- decided-by [check-quality-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/check-quality-topic.md)
- decided-by [check-qa-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/check-qa-topic.md)
- decided-by [check-model-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/check-model-topic.md)
- decided-by [critic-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/critic-topic.md)
- decided-by [fixer-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/fixer-topic.md)
- decided-by [ci-fix-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/ci-fix-topic.md)
- decided-by [gate-fix-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/gate-fix-topic.md)
- decided-by [conflict-fix-topic](../capability/rhiza/workflows/built-in-workflows/built-in-review/review-guide/conflict-fix-topic.md)
- decided-by [shared-severities-block](../capability/rhiza/workflows/built-in-workflows/plan-workflow/shared-severities-block.md)
- decided-by [planner-topic](../capability/rhiza/workflows/built-in-workflows/plan-workflow/planner-topic.md)
- decided-by [plan-reviewer-topic](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviewer-topic.md)
- decided-by [plan-critic-topic](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-critic-topic.md)
- decided-by [plan-reviser-topic](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-reviser-topic.md)
- decided-by [plan-role-skills-installed](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-role-skills-installed.md)
- decided-by [plan-prompts-one-line](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-prompts-one-line.md)

## History

- `01M76SV05752A64QJNN4QTD2KH` Efficient build cycle. Owner, 2026-11-12: "I'm disappointed with the overall performance of a full cycle of build runs. It takes a long time, and seems to be error prone. ... I like this process overall but it's just not efficient." On the shipit/forge prompts: "I also think we may have lost something by not including the shipit/forge prompts in the workflow steps." "Those should ship as the default build tools, and we also want to use them as our own." After the orchestrator's summary of both analyses and its proposal of one milestone ("Efficient build cycle", afk prompts as shipped defaults plus the four guards first, then the CI-round cuts and the error sources behind held runs, with the owner deciding on parallel leaves, a cheaper verifier model and dropping the critic), the owner said: "Sounds good". Evidence to start from (read it first; numbers are measured unless marked estimated): ~/orchestrator/perf-2026-11-12/report.md (233 run journals and 699 CI runs: where a ship-change's time goes, unit costs per failure cause, ranked recommendations, statistics), prompt-gap.md (afk at github.com/dphaener/afk plugins/afk, the merged shipit/forge prompts, compared role by role with Rhiza's step prompts, measured afk archives: median 39m per task vs Rhiza 29m plan + 121m ship, rework 12% vs 31% build retries and 7% vs 58% review fix cycles; 15 ranked recommendations with estimated effect, and a 12 change set sketch), tables.md and tables-recent.md. The analysis scripts are in the orchestrator's scratchpad (analyze.py, rerun after changes; a resource sampler is collecting the next 30 runs). Already filed and not to be duplicated: model-only CI for plan PRs (01M764CJ), pre-push self-check in every run (01M76R1E), non-test CI failures routed to ci-fix (01M76A58), nothing uncommitted on main (01M76CSJ), log-entry write guard and repair (01M76Q44), plan review rounds that fail to record (01M76DWZ). Related milestones: 32 (workflow cleanup and language-neutral defaults) and the agent docs milestone (agents/ with orchestrator.md).
