# delay-repro-plan-interview

source · `01M6ZCKVQ9GM1J8TG067QZT57X`

## quotes

```json
[]
```

## raw_sha256

316971058618220be8d3642b50dee3799c746d8b5b576a68bed6df46cb164d38

## session

d5c3217c-51fd-4fba-9e9f-01a1bec67441

## summary

Autonomous plan run for change set 01M6ZCC (Builders and fixers reproduce). The intent: two nebari runs on 2026-11-09 reproduced a flaky test by starting 8 and 14 busy loops on the shared 4-core VM, holding load at 16-22 for 15-30+ minutes and slowing every other live run, while load-based repro barely worked (45/45 passed) and a 500 ms injected delay on the suspected request reproduced 10/10. No owner was present; the planner asked three questions, which the arbiter decided in decision-delay-repro-scope: the change set serves milestone-25 (empty backlog: bugs and polish), like the earlier unstable-CI guidance and flaky-test fixes; it is guidance only, putting the rule 'reproduce timing flakes by injecting a delay at the suspected race, never by busy loops or parallel stress on the host' into the build and review guide topics (the source of the rhiza-build and rhiza-review skills) and into this repo's CLAUDE.md beside the unstable-CI reproduce-locally rule; the optional watchdog flagging runs' descendant processes at 100% CPU is left out, to be planned as its own change set later. Planned: one leaf, delay-repro-guidance, under self-check next to unstable-ci-guidance, checked by the guide text tests and the CLAUDE.md test.

## title

Plan interview: builders and fixers reproduce flakes by delay, not CPU load

## Links

### In

- derived-from [decision-delay-repro-scope](../decision/decision-delay-repro-scope.md)

## History

- `01M6ZCC27HQ7NMBTZ74D1CJFBV` Builders and fixers reproduce flaky tests by generating CPU load on the shared host. On 2026-11-09 two nebari runs did it: 01M6Z1RG ran 8 'bash' busy loops, and 01M6Z5P4 (session 42a0d9da) ran 14 'sh -c while :; do :; done' loops. Each held the 4-core VM at load 16–22 for 15–30+ min and slowed every other live run; I stopped them by PID. Load-based repro also barely works: 45/45 passed under load, while a 500 ms injected delay on the suspected request reproduced 10/10 (nebari PR #90). The build/fix guidance (rhiza guide build, the rhiza-build and rhiza-review skills) should say: reproduce timing flakes by injecting a delay at the suspected race, never by busy loops or parallel stress on the host. Optionally the watchdog flags a run's descendant processes that sit at 100% CPU for minutes.
