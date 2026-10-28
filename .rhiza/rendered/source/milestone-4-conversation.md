# milestone-4-conversation

source · `01M61R316Z5KYNPHDHFSB369SD`

## quotes

```json
[]
```

## raw_sha256

1e247df640bd1299c1b2e4630c25f768978bf163a4bb0e60d8db19d4f6b9024d

## session

8cd68bc7-9786-4df4-8723-01a1831302a8

## summary

The owner and Rhiza's plan workflow (a Claude Code planner in Rex) planned milestone 4 on 2026-10-28, after milestone 3 shipped Rhiza's plan, build, review and ship-change workflows, with a Claude Code orchestrator session still choosing and watching runs by hand. Every answer was the owner picking the planner's recommended option, so all wording is the planner's, accepted by the owner; the owner typed nothing. Decided: milestone 4 is orchestration (layer 5), so milestone 5 is driven by Rhiza with no orchestrator session, and milestone 3's gaps come first. The orchestrator is a deterministic daemon in the same binary (`rhiza daemon`), calling an arbiter only in autonomous mode. It starts ship-change runs by itself on accepted, ready change sets that serve a milestone and have leaves, in roadmap order, under a concurrency cap (default 2). It runs in the foreground of a Rex tab, one per repo through a lock under the git common dir, polls the roadmap and run journals with no socket, starts each run in its own Rex tab, passes --autonomous through, and reports through `rhiza daemon status`. It resumes an interrupted run once and leaves failed runs to the inbox. The owner deals with it through an inbox (`rhiza inbox`): questions, sign-offs and failed runs across runs, each question with a recommendation from an advisor session that reads slice and why and never answers. `rhiza inbox open` focuses the asking session's tab, where questions are answered; sign-offs are acted on from the inbox; there is a Rex notification for each new question, sign-off and failed run. The usage governor is a concurrency cap plus a pause on the Claude usage limit. Done points: milestone 5's first change sets shipped with the daemon orchestrating; runs that finish without hand repair; make check still green, reads within budget and no socket in the default tests. All 15 draft gap change sets and the proposed questions-with-recommendation change set (the inbox) serve milestone 4. Two new change sets: the daemon and a usage-limit pause. Out: a TUI or web UI, unattended machines, token budgets, plan runs started by the daemon, conflict detection across change sets, and plan review against runs in flight. Planned now: one live run per change set (a run is live until its journal finishes, and a new `rhiza run abandon` ends a stuck run), the daemon and the inbox; the rest keep their intent only.

## title

Milestone 4 conversation

## Links

### In

- derived-from [answer-8cd68bc7-1](../decision/answer-8cd68bc7-1.md)
- derived-from [answer-8cd68bc7-2](../decision/answer-8cd68bc7-2.md)
- derived-from [answer-8cd68bc7-3](../decision/answer-8cd68bc7-3.md)
- derived-from [answer-8cd68bc7-4](../decision/answer-8cd68bc7-4.md)
- derived-from [answer-8cd68bc7-5](../decision/answer-8cd68bc7-5.md)
- derived-from [answer-8cd68bc7-6](../decision/answer-8cd68bc7-6.md)
- derived-from [answer-8cd68bc7-7](../decision/answer-8cd68bc7-7.md)
- derived-from [answer-8cd68bc7-8](../decision/answer-8cd68bc7-8.md)
- derived-from [answer-8cd68bc7-9](../decision/answer-8cd68bc7-9.md)
- derived-from [answer-8cd68bc7-10](../decision/answer-8cd68bc7-10.md)
- derived-from [answer-8cd68bc7-11](../decision/answer-8cd68bc7-11.md)
- derived-from [answer-8cd68bc7-12](../decision/answer-8cd68bc7-12.md)
- derived-from [answer-8cd68bc7-13](../decision/answer-8cd68bc7-13.md)
- derived-from [answer-8cd68bc7-14](../decision/answer-8cd68bc7-14.md)
- derived-from [answer-8cd68bc7-15](../decision/answer-8cd68bc7-15.md)
- derived-from [answer-8cd68bc7-16](../decision/answer-8cd68bc7-16.md)
- derived-from [answer-8cd68bc7-17](../decision/answer-8cd68bc7-17.md)
- derived-from [answer-8cd68bc7-18](../decision/answer-8cd68bc7-18.md)
- derived-from [milestone-4](../intent/milestone-4.md)

## History

- `01M61H5VY8MD7E7SVS0HQQHQJ0` Plan milestone 4 with the owner: its intent, its done points and the change sets that serve it, planned by Rhiza's plan workflow.
