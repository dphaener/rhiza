# question-d5c3217c-1

question · `01M6ZCH61HJRA2WPADKS4GY4BM`

## session

d5c3217c-51fd-4fba-9e9f-01a1bec67441

## statement

Change set 01M6ZCC (Builders and fixers reproduce) serves no intent yet, and milestone-30 (ready to share) is the only open milestone on the roadmap. Recent guidance and flake fixes (#403 Stop using CI as a debugger, #405 Fix flaky test on main) served milestone-25 (empty backlog: bugs and polish). What should this change set serve? / The intent's watchdog part is optional: the daemon flags a run's descendant processes that sit at 100% CPU for minutes. How should change set 01M6ZCC handle it? / Should the rule 'reproduce timing flakes by injecting a delay at the suspected race, never by busy loops or parallel stress on the host' go only into the build and review guide topics (the source of the rhiza-build and rhiza-review skills), or also into this repo's CLAUDE.md?

## Links

### In

- answers [decision-delay-repro-scope](../decision/decision-delay-repro-scope.md)

## History

- `01M6ZCC27HQ7NMBTZ74D1CJFBV` Builders and fixers reproduce flaky tests by generating CPU load on the shared host. On 2026-11-09 two nebari runs did it: 01M6Z1RG ran 8 'bash' busy loops, and 01M6Z5P4 (session 42a0d9da) ran 14 'sh -c while :; do :; done' loops. Each held the 4-core VM at load 16–22 for 15–30+ min and slowed every other live run; I stopped them by PID. Load-based repro also barely works: 45/45 passed under load, while a 500 ms injected delay on the suspected request reproduced 10/10 (nebari PR #90). The build/fix guidance (rhiza guide build, the rhiza-build and rhiza-review skills) should say: reproduce timing flakes by injecting a delay at the suspected race, never by busy loops or parallel stress on the host. Optionally the watchdog flags a run's descendant processes that sit at 100% CPU for minutes.
