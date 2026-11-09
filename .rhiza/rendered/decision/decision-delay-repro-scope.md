# decision-delay-repro-scope

decision · `01M6ZCJWW14QZ3FM14X8EGG7V0`

## reasoning

Milestone-30 is about making Rhiza ready to share (install, README, docs, reference, public); agent build/fix discipline does not serve any of its parts. The closest precedents, #403 (Stop using CI as a debugger) and #405 (flaky test on main), served milestone-25, which explicitly covers bugs and polish of the daemon's autonomous building, and this rule is the same kind of guidance. #403 put its rule in both CLAUDE.md and the guides; CLAUDE.md already holds the reproduce-locally rule this one refines, and agents in this repo read CLAUDE.md first, so placing it there too keeps the two rules together and costs one paragraph. The watchdog is marked optional in the intent and is daemon code (process sampling, thresholds, settings) with its own design questions; bundling it would turn a small guidance change into a code change and delay the part that stops the harm now.

## statement

Change set 01M6ZCC serves milestone-25 (empty backlog: bugs and polish), not milestone-30. It is guidance only: the delay-injection rule ('reproduce timing flakes by injecting a delay at the suspected race, never by busy loops or parallel stress on the host') goes into the build and review guide topics (and so the rhiza-build and rhiza-review skills) and also into this repo's CLAUDE.md, beside the unstable-CI reproduce-locally rule. The optional watchdog for runs' descendant processes at 100% CPU is left out of this change set; it can be planned as its own change set later.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-d5c3217c-1](../question/question-d5c3217c-1.md)
- derived-from [delay-repro-plan-interview](../source/delay-repro-plan-interview.md)

### In

- decided-by [delay-repro-guidance](../capability/rhiza/integrity/self-check/delay-repro-guidance.md)

## History

- `01M6ZCC27HQ7NMBTZ74D1CJFBV` Builders and fixers reproduce flaky tests by generating CPU load on the shared host. On 2026-11-09 two nebari runs did it: 01M6Z1RG ran 8 'bash' busy loops, and 01M6Z5P4 (session 42a0d9da) ran 14 'sh -c while :; do :; done' loops. Each held the 4-core VM at load 16–22 for 15–30+ min and slowed every other live run; I stopped them by PID. Load-based repro also barely works: 45/45 passed under load, while a 500 ms injected delay on the suspected request reproduced 10/10 (nebari PR #90). The build/fix guidance (rhiza guide build, the rhiza-build and rhiza-review skills) should say: reproduce timing flakes by injecting a delay at the suspected race, never by busy loops or parallel stress on the host. Optionally the watchdog flags a run's descendant processes that sit at 100% CPU for minutes.
