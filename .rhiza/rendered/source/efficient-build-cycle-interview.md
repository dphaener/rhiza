# efficient-build-cycle-interview

source · `01M76XM73SV2HZ5FRGC8CMX965`

## quotes

```json
[
  {
    "by": "owner",
    "text": "I'm disappointed with the overall performance of a full cycle of build runs. It takes a long time, and seems to be error prone. ... I like this process overall but it's just not efficient.",
    "where": "2026-11-12, request to the orchestrator"
  },
  {
    "by": "owner",
    "text": "I also think we may have lost something by not including the shipit/forge prompts in the workflow steps.",
    "where": "2026-11-12, request to the orchestrator"
  },
  {
    "by": "owner",
    "text": "Those should ship as the default build tools, and we also want to use them as our own.",
    "where": "2026-11-12, request to the orchestrator"
  },
  {
    "by": "owner",
    "text": "Sounds good",
    "where": "2026-11-12, on the orchestrator's milestone proposal"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Guards: Done needs a record, Pre-handoff checklist, Command budget",
    "where": "interview"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Role topics, 3 change sets; afk's content, Rhiza's words",
    "where": "interview"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Critic: keep, with downgrade; verifier: yes, tiered by risk; parallel leaves: later, after leaf files",
    "where": "interview"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Post-land CI: model-only since last green; merge rounds: yes, same rule",
    "where": "interview"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Targets: rework halved, ship <=40m",
    "where": "interview"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Extras: leaf contract, verbatim retry + run notes, model lints in review; queue journaling: yes",
    "where": "interview"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Plan now: the three guards, CI since last green, held-run error sources, critic downgrade; scope check: yes, fold it",
    "where": "interview"
  }
]
```

## raw_sha256

8dcaaefb604d95fb7559da994183cce51e88cd9108a8676347c0dc3fc29d74f6

## session

e3a4887c-f758-434a-a9a4-01a1cd9de36f

## summary

Planned milestone 34, Efficient build cycle, from the 2026-11-12 performance report (233 run journals, 699 CI runs) and the afk prompt-gap analysis. The owner had asked for afk's shipit/forge prompts to ship as Rhiza's default build tools and accepted the orchestrator's milestone proposal. In this interview the owner picked: guards are done-needs-a-record, a builder/fixer pre-handoff checklist and a command budget (the .rhiza/log tool guard and generate-before-push are already filed as 01M76QJW and 01M76RH0; a git-safety deny was not picked); afk's prompts become one guide topic per role in three change sets (build, review, plan roles), keeping afk's content in Rhiza's words; the critic stays but can downgrade; low-risk leaves get a cheaper verifier tiered by risk; parallel leaves come later, after leaves name their files; the post-land CI round is cut by running only model checks when nothing outside .rhiza/ changed since the PR's last green commit, and the same rule covers merges of main; targets are build retries <=15% of leaves, second review cycles <=20% of runs, ship-change median <=40m and mean held time <=10m over the 30 runs after; the leaf contract, verbatim retry with run notes, deterministic model lints and queue-wait journaling are in, a lean session start is not; the scope check folds into compliance. Planned now: the three guards, CI since last green, held-run error sources and critic downgrade. Open: rhiza ci install writes only a rhiza check workflow, so where the since-last-green rule ships for other projects is parked as a question.

## title

Efficient build cycle milestone, with the owner

## Links

### In

- derived-from [afk-prompts-shipped-default](../decision/afk-prompts-shipped-default.md)
- derived-from [efficient-cycle-milestone-shape](../decision/efficient-cycle-milestone-shape.md)
- derived-from [answer-e3a4887c-1](../decision/answer-e3a4887c-1.md)
- derived-from [answer-e3a4887c-2](../decision/answer-e3a4887c-2.md)
- derived-from [answer-e3a4887c-3](../decision/answer-e3a4887c-3.md)
- derived-from [answer-e3a4887c-4](../decision/answer-e3a4887c-4.md)
- derived-from [answer-e3a4887c-5](../decision/answer-e3a4887c-5.md)
- derived-from [answer-e3a4887c-6](../decision/answer-e3a4887c-6.md)
- derived-from [answer-e3a4887c-7](../decision/answer-e3a4887c-7.md)
- derived-from [answer-e3a4887c-8](../decision/answer-e3a4887c-8.md)
- derived-from [answer-e3a4887c-9](../decision/answer-e3a4887c-9.md)
- derived-from [answer-e3a4887c-10](../decision/answer-e3a4887c-10.md)
- derived-from [answer-e3a4887c-11](../decision/answer-e3a4887c-11.md)
- derived-from [answer-e3a4887c-12](../decision/answer-e3a4887c-12.md)
- derived-from [answer-e3a4887c-13](../decision/answer-e3a4887c-13.md)

## History

- `01M76SV05752A64QJNN4QTD2KH` Efficient build cycle. Owner, 2026-11-12: "I'm disappointed with the overall performance of a full cycle of build runs. It takes a long time, and seems to be error prone. ... I like this process overall but it's just not efficient." On the shipit/forge prompts: "I also think we may have lost something by not including the shipit/forge prompts in the workflow steps." "Those should ship as the default build tools, and we also want to use them as our own." After the orchestrator's summary of both analyses and its proposal of one milestone ("Efficient build cycle", afk prompts as shipped defaults plus the four guards first, then the CI-round cuts and the error sources behind held runs, with the owner deciding on parallel leaves, a cheaper verifier model and dropping the critic), the owner said: "Sounds good". Evidence to start from (read it first; numbers are measured unless marked estimated): ~/orchestrator/perf-2026-11-12/report.md (233 run journals and 699 CI runs: where a ship-change's time goes, unit costs per failure cause, ranked recommendations, statistics), prompt-gap.md (afk at github.com/dphaener/afk plugins/afk, the merged shipit/forge prompts, compared role by role with Rhiza's step prompts, measured afk archives: median 39m per task vs Rhiza 29m plan + 121m ship, rework 12% vs 31% build retries and 7% vs 58% review fix cycles; 15 ranked recommendations with estimated effect, and a 12 change set sketch), tables.md and tables-recent.md. The analysis scripts are in the orchestrator's scratchpad (analyze.py, rerun after changes; a resource sampler is collecting the next 30 runs). Already filed and not to be duplicated: model-only CI for plan PRs (01M764CJ), pre-push self-check in every run (01M76R1E), non-test CI failures routed to ci-fix (01M76A58), nothing uncommitted on main (01M76CSJ), log-entry write guard and repair (01M76Q44), plan review rounds that fail to record (01M76DWZ). Related milestones: 32 (workflow cleanup and language-neutral defaults) and the agent docs milestone (agents/ with orchestrator.md).
