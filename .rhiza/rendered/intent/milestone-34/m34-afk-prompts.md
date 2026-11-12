# m34-afk-prompts

intent · `01M76XMNBAT0FWMMA5G09W751S` · parent [milestone-34](../milestone-34.01M76XMNA9ATHT77CB9GCH44P3.md)

## statement

Every run role's shipped prompt carries afk's shipit/forge rules and checklists in Rhiza's words, one language-neutral guide topic per role on afk's six-section skeleton with shared blocks written once, and Rhiza's own repo runs them.

## Links

### Out

- decided-by [afk-prompts-shipped-default](../../decision/afk-prompts-shipped-default.md)

### In

- serves change set `01M76XR5Y1N5E699AT2WMJMQJC` (Build roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton (task, expected outcome, read first, must do, must not do, context) with shared blocks written once: builder, verifier (re-run reported checks and acceptance, findings a fixer can act on) and builder retry; the builtin build prompts shrink to one line. Source: prompt-gap.md sections 2.3, 2.4 and 5 (CS1).)
- serves change set `01M76XRETZFTD2TQ5GS3WYAERA` (Review roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: compliance (absorbing the scope check: under-built, scope creep, unaccounted files), quality (impact radius, checklist), qa (clean state, scratch build outside PATH, happy and failure cases), model, critic, fixer, ci-fix, gate-fix and conflict-fix; checks are isolated from build history; the scope check leaves review.toml. Source: prompt-gap.md sections 2.5-2.8 and 5 (CS1, CS8).)
- serves change set `01M76XRPN6WYMZSV2GWA05Y1FM` (Plan roles get afk's shipit/forge prompts in Rhiza's words, one guide topic each on the six-section skeleton: planner (a codebase research pass before deciding, a self-check re-reading the slice after writing), plan reviewer (Crucible checklist, uncovered-intent and oversized-leaf kinds, nothing a deterministic check already checks), plan critic and reviser. Source: prompt-gap.md sections 2.1, 2.2 and 5 (CS1, rec 15).)

## History

- `01M76SV05752A64QJNN4QTD2KH` Efficient build cycle. Owner, 2026-11-12: "I'm disappointed with the overall performance of a full cycle of build runs. It takes a long time, and seems to be error prone. ... I like this process overall but it's just not efficient." On the shipit/forge prompts: "I also think we may have lost something by not including the shipit/forge prompts in the workflow steps." "Those should ship as the default build tools, and we also want to use them as our own." After the orchestrator's summary of both analyses and its proposal of one milestone ("Efficient build cycle", afk prompts as shipped defaults plus the four guards first, then the CI-round cuts and the error sources behind held runs, with the owner deciding on parallel leaves, a cheaper verifier model and dropping the critic), the owner said: "Sounds good". Evidence to start from (read it first; numbers are measured unless marked estimated): ~/orchestrator/perf-2026-11-12/report.md (233 run journals and 699 CI runs: where a ship-change's time goes, unit costs per failure cause, ranked recommendations, statistics), prompt-gap.md (afk at github.com/dphaener/afk plugins/afk, the merged shipit/forge prompts, compared role by role with Rhiza's step prompts, measured afk archives: median 39m per task vs Rhiza 29m plan + 121m ship, rework 12% vs 31% build retries and 7% vs 58% review fix cycles; 15 ranked recommendations with estimated effect, and a 12 change set sketch), tables.md and tables-recent.md. The analysis scripts are in the orchestrator's scratchpad (analyze.py, rerun after changes; a resource sampler is collecting the next 30 runs). Already filed and not to be duplicated: model-only CI for plan PRs (01M764CJ), pre-push self-check in every run (01M76R1E), non-test CI failures routed to ci-fix (01M76A58), nothing uncommitted on main (01M76CSJ), log-entry write guard and repair (01M76Q44), plan review rounds that fail to record (01M76DWZ). Related milestones: 32 (workflow cleanup and language-neutral defaults) and the agent docs milestone (agents/ with orchestrator.md).
