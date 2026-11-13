# spawn-handback-depend-on-contract-writers

decision · `01M7AMFBEDG9C7Q7FRZV0Z0BZW`

## reasoning

spawn-reuses-contract-and-waits-on-holds already decided that this change set reuses the contract change set's writers rather than building its own, and both leaves' statements and acceptance say so: run-spawn extends 'the contract change set's brief writer' with a spawn mode, run-handback passes its cause to 'the contract change set's handback.md writer' and sends 'the handback text the contract change set's agent-invoke-template defines'. run-brief-files and agent-invoke-template are created by 01M79KGHYF and are not yet landed, yet rhiza show lists no depends-on from either leaf to them, so the build order (which counts depends-on in or out of the change set) could start these leaves first and the builder would find nothing to extend. Reversing the Decision instead (building separate writers here) would duplicate code and let the two drift, which that Decision and conductor-milestone-direction's few-primitives direction rejected. Adding the three depends-on links is the smallest fix that matches the finding, the critic's verdict and the existing Decision.

## statement

Settle it by recording the dependency in the leaf graph: run-spawn depends-on run-brief-files, and run-handback depends-on run-brief-files and agent-invoke-template (in addition to its existing handback-route and gate-step). The Decision spawn-reuses-contract-and-waits-on-holds stands as written; the links make build order honour it, so neither leaf is built before the contract change set's brief.md/handback.md writer and handback text exist.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7ah6q9g3ytjk9antmptwxvm-r3-f1](../question/review-01m7ah6q9g3ytjk9antmptwxvm-r3-f1.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
