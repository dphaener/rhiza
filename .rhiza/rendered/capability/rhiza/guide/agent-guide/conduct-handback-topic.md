# conduct-handback-topic

capability · `01M7MJXKVHQXE1B5MWRMZNM6RD` · parent [agent-guide](../agent-guide.md)

## acceptance

go test -trimpath ./internal/guide -run 'TestConductHandback|TestConduct|TestList' shows the conduct-handback topic listed and under TopicCap, the conduct topic's Read first naming rhiza guide conduct-handback and staying under TopicCap, and conduct-handback naming handback.md's section titles as internal/run/brief.go renders them (Failing step, Outcome, Failing commands, Output tail, Cause, Failed CI runs, Conflicted files), rhiza run build-check, rhiza run review-check, rhiza check, rhiza ship gate, rhiza ship open-pr, rhiza ship rebase --no-push, the roles verifier, critic, fixer, gate-fix, ci-fix and conflict-fix, handback_cap and rhiza session done; the test parses the built-in ship-change with the workflow package and requires the topic to name every step whose on.* target is handback:conduct and every command of its gate step; (failure case) the test fails naming the step or command when ship-change hands back from a step, or its gate runs a command, the topic does not name, and the topic's Routes section (one entry per hand-back step) naming a step whose on.* targets include no handback:conduct (main-ci, close-out, ship, deploy-check, derived from the parsed ship-change) as a route fails it naming that step; the wait-ci route naming the ci-fix role fails it; go test -trimpath ./cmd/rhiza -run TestScripts/guide passes.

## files

internal/guide/text/conduct-handback.md, internal/guide/conduct_handback_test.go, internal/guide/guide.go, internal/guide/guide_test.go, internal/guide/text/conduct.md, internal/guide/conduct_test.go, cmd/rhiza/testdata/script/guide.txtar

## must_not

Change the hand-back mechanism, handback.md's rendering, ship-change.toml or any fixer role topic: guide text and its tests only; grow the conduct topic past TopicCap; tell the conductor to fix anything itself

## statement

The guide has a conduct-handback topic, named in the conduct topic's Read first, giving the shape of handback.md (Failing step, Outcome, Failing commands with output paths, Output tail, Cause, Failed CI runs, Conflicted files) and, for each step of the conductor default ship-change that hands back to conduct, the route a conductor takes as conduct-handback-routes decides: the gate's four commands, signoff changes, push, wait-ci's timeout and merge (CI red, conflicting PR, moved main via rhiza ship rebase --no-push), each by spawning the named role, then land, render and rhiza session done; plus handback_cap and that the steps after the merge never hand back.

## Links

### Out

- decided-by [conduct-guide-three-topics](../../../../decision/conduct-guide-three-topics.md)
- decided-by [conduct-handback-routes](../../../../decision/conduct-handback-routes.md)

### In

- verifies [conduct-handback-topic-verified](../../../../evidence/conduct-handback-topic-verified.md)
- verifies [conduct-handback-topic-qa-1](../../../../evidence/conduct-handback-topic-qa-1.md)
- verifies [conduct-handback-topic-qa-2](../../../../evidence/conduct-handback-topic-qa-2.md)

## History

- `01M7MJN83Y6EF727C2T0AF3Y5D` The rhiza-conduct skill and its guide topics are enough on their own: a conductor outside this repo has no Rhiza source to read, so everything a conductor looked up in the code here (step roles, handback routes, workflow-new defaults, check roles) is in the guide
