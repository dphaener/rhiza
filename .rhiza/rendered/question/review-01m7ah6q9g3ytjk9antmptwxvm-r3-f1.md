# review-01m7ah6q9g3ytjk9antmptwxvm-r3-f1

question · `01M7AMCGDYJ1ATX2XK77ZJ9AN0`

## statement

Plan review round 3 left finding f1 (contradiction on spawn-reuses-contract-and-waits-on-holds): The Decision says this change set depends on the contract change set (01M79KGHYF) and calls its writers, and run-spawn and run-handback extend run-brief-files' brief.md/handback.md writer and send agent-invoke-template's handback text, but no depends-on link records it: rhiza change show lists run-spawn with needs [] and run-handback with needs [gate-step, handback-route] only, and neither run-brief-files nor agent-invoke-template is linked. Build order (leaf needs count depends-on in or out of the change set) therefore lets run-spawn and run-handback be built before the writers they extend exist. Add depends-on from run-spawn to run-brief-files and from run-handback to run-brief-files and agent-invoke-template. The critic confirmed it: Verified: the Decision spawn-reuses-contract-and-waits-on-holds says this change set depends on the contract change set 01M79KGHYF and calls its writers, and run-spawn's touches name 'the contract change set's brief.md writer (a spawn mode...)'. Yet rhiza show run-spawn lists no depends-on at all, and run-handback depends only on handback-route and gate-step; neither links run-brief-files or agent-invoke-template, and rhiza check does not report it. Build order can therefore start run-spawn/run-handback before the writers they extend exist. A real contradiction between the Decision and the leaf graph; major stands. How should it be settled?

## Links

### In

- answers [spawn-handback-depend-on-contract-writers](../decision/spawn-handback-depend-on-contract-writers.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
