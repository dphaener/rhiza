# efficient-cycle-milestone-shape

decision · `01M76XMNNPW27705N7G6XFBDCH`

## quote

Sounds good

## statement

Milestone 34 puts afk's prompts as shipped defaults and the guards first, then the CI-round cuts and the error sources behind held runs; the critic stays and may downgrade a finding, the verifier's model is tiered by leaf risk, and parallel leaves come later, after leaves name their files

## Links

### Out

- derived-from [efficient-build-cycle-interview](../source/efficient-build-cycle-interview.md)

### In

- decided-by [milestone-35](../intent/milestone-35.md)
- decided-by [done-refuses-unrecorded-verifier](../capability/rhiza/workflows/workflow-runs/agent-step-done-signal/done-refuses-unrecorded-verifier.md)
- decided-by [done-refuses-unrecorded-check](../capability/rhiza/workflows/workflow-runs/agent-step-done-signal/done-refuses-unrecorded-check.md)
- decided-by [done-own-command-guide](../capability/rhiza/guide/done-own-command-guide.md)
- decided-by [command-budget-guide](../capability/rhiza/guide/command-budget-guide.md)

## History

- `01M76SV05752A64QJNN4QTD2KH` Efficient build cycle. Owner, 2026-11-12: "I'm disappointed with the overall performance of a full cycle of build runs. It takes a long time, and seems to be error prone. ... I like this process overall but it's just not efficient." On the shipit/forge prompts: "I also think we may have lost something by not including the shipit/forge prompts in the workflow steps." "Those should ship as the default build tools, and we also want to use them as our own." After the orchestrator's summary of both analyses and its proposal of one milestone ("Efficient build cycle", afk prompts as shipped defaults plus the four guards first, then the CI-round cuts and the error sources behind held runs, with the owner deciding on parallel leaves, a cheaper verifier model and dropping the critic), the owner said: "Sounds good". Evidence to start from (read it first; numbers are measured unless marked estimated): ~/orchestrator/perf-2026-11-12/report.md (233 run journals and 699 CI runs: where a ship-change's time goes, unit costs per failure cause, ranked recommendations, statistics), prompt-gap.md (afk at github.com/dphaener/afk plugins/afk, the merged shipit/forge prompts, compared role by role with Rhiza's step prompts, measured afk archives: median 39m per task vs Rhiza 29m plan + 121m ship, rework 12% vs 31% build retries and 7% vs 58% review fix cycles; 15 ranked recommendations with estimated effect, and a 12 change set sketch), tables.md and tables-recent.md. The analysis scripts are in the orchestrator's scratchpad (analyze.py, rerun after changes; a resource sampler is collecting the next 30 runs). Already filed and not to be duplicated: model-only CI for plan PRs (01M764CJ), pre-push self-check in every run (01M76R1E), non-test CI failures routed to ci-fix (01M76A58), nothing uncommitted on main (01M76CSJ), log-entry write guard and repair (01M76Q44), plan review rounds that fail to record (01M76DWZ). Related milestones: 32 (workflow cleanup and language-neutral defaults) and the agent docs milestone (agents/ with orchestrator.md).
- `01M76YXNTKWHGGAHF7MNH4G7GF` apply -
