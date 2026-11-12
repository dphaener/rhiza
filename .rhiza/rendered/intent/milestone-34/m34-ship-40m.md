# m34-ship-40m

intent · `01M76XMNGHCD7XX82ZCSRPRGPC` · parent [milestone-34](../milestone-34.01M76XMNA9ATHT77CB9GCH44P3.md)

## statement

Median ship-change wall time is at most 40 minutes over the 30 runs after the milestone (baseline Nov 9-12: 55 minutes), measured by the perf analysis and recorded as Evidence, and the daemon journals why an accepted change set waits before its run starts.

## Links

### In

- serves change set `01M76XS745J8S0J74SZFEHQ17Z` (Build verifiers run on a cheaper model for low-risk leaves and on the strongest model for risky ones (test infrastructure, hooks, CI, workflow files, migrations, retries), still writing non-author Evidence. Evidence: the verifier is the largest single bucket, 16.9% of ship-change wall (prompt-gap rec 13).)
- serves change set `01M76XSY0YM17EQ9K24CC4JRAK` (The daemon journals why an accepted change set waits before its run starts (no free slot, load hold, usage hold, dependency, policy), and the perf analysis can read it. Evidence: accepted to ship-change start is 27m RECENT and 60m ALL median, cause not journaled (perf report rec 10).)
- serves change set `01M76XTN06GSM9Z3V7CHKHFX7Q` (Leaves whose files are disjoint build in parallel within a run, with a lock per file and bounded by the daemon's load hold. Evidence: build grows 10-20m per leaf; 29 of 119 runs had 4 or more leaves (prompt-gap rec 12). Planned after the leaf contract, which gives leaves their files.)

## History

- `01M76SV05752A64QJNN4QTD2KH` Efficient build cycle. Owner, 2026-11-12: "I'm disappointed with the overall performance of a full cycle of build runs. It takes a long time, and seems to be error prone. ... I like this process overall but it's just not efficient." On the shipit/forge prompts: "I also think we may have lost something by not including the shipit/forge prompts in the workflow steps." "Those should ship as the default build tools, and we also want to use them as our own." After the orchestrator's summary of both analyses and its proposal of one milestone ("Efficient build cycle", afk prompts as shipped defaults plus the four guards first, then the CI-round cuts and the error sources behind held runs, with the owner deciding on parallel leaves, a cheaper verifier model and dropping the critic), the owner said: "Sounds good". Evidence to start from (read it first; numbers are measured unless marked estimated): ~/orchestrator/perf-2026-11-12/report.md (233 run journals and 699 CI runs: where a ship-change's time goes, unit costs per failure cause, ranked recommendations, statistics), prompt-gap.md (afk at github.com/dphaener/afk plugins/afk, the merged shipit/forge prompts, compared role by role with Rhiza's step prompts, measured afk archives: median 39m per task vs Rhiza 29m plan + 121m ship, rework 12% vs 31% build retries and 7% vs 58% review fix cycles; 15 ranked recommendations with estimated effect, and a 12 change set sketch), tables.md and tables-recent.md. The analysis scripts are in the orchestrator's scratchpad (analyze.py, rerun after changes; a resource sampler is collecting the next 30 runs). Already filed and not to be duplicated: model-only CI for plan PRs (01M764CJ), pre-push self-check in every run (01M76R1E), non-test CI failures routed to ci-fix (01M76A58), nothing uncommitted on main (01M76CSJ), log-entry write guard and repair (01M76Q44), plan review rounds that fail to record (01M76DWZ). Related milestones: 32 (workflow cleanup and language-neutral defaults) and the agent docs milestone (agents/ with orchestrator.md).
