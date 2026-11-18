# conduct-handback-topic

capability · `01M7MJXKVHQXE1B5MWRMZNM6RD` · parent [agent-guide](../agent-guide.md)

## acceptance

go test -trimpath ./internal/guide -run 'TestConductHandback|TestConduct|TestList' shows the conduct-handback topic listed and under TopicCap, the conduct topic's Read first naming rhiza guide conduct-handback and staying under TopicCap, and conduct-handback naming handback.md's section titles as internal/run/brief.go renders them (Failing step, Outcome, Failing commands, Output tail, Cause, Failed CI runs, Conflicted files), rhiza run build-check, rhiza run review-check, rhiza check, rhiza ship gate, rhiza ship open-pr, the roles verifier, critic, fixer, gate-fix, ci-fix and conflict-fix, handback_cap and rhiza session done, with routes for pr-green and rebase and none for merge, and not naming rhiza ship rebase --no-push; the test parses the built-in ship-change with the workflow package and requires the topic to name every step whose on.* target is handback:conduct and every command of its gate step; (failure case) the test fails naming the step or command when ship-change hands back from a step, or its gate runs a command, the topic does not name, and the topic's Routes section (one entry per hand-back step) naming a step whose on.* targets include no handback:conduct (merge, main-ci, close-out, ship, deploy-check, derived from the parsed ship-change) as a route fails it naming that step; the wait-ci route naming the ci-fix role fails it; a route telling the conductor to run rhiza ship rebase, with or without --no-push, fails it; go test -trimpath ./cmd/rhiza -run TestScripts/guide passes.

## files

internal/guide/text/conduct-handback.md, internal/guide/conduct_handback_test.go, internal/guide/guide.go, internal/guide/guide_test.go, internal/guide/text/conduct.md, internal/guide/conduct_test.go, cmd/rhiza/testdata/script/guide.txtar

## must_not

Change the hand-back mechanism, handback.md's rendering, ship-change.toml or any fixer role topic: guide text and its tests only; grow the conduct topic past TopicCap; tell the conductor to fix anything itself

## statement

The guide has a conduct-handback topic, named in the conduct topic's Read first, giving the shape of handback.md (Failing step, Outcome, Failing commands with output paths, Output tail, Cause, Failed CI runs, Conflicted files) and, for each step of the conductor default ship-change that hands back to conduct, the route a conductor takes as conduct-handback-routes-ship-rebase decides: the gate's four commands, signoff changes, push, wait-ci's timeout, pr-green (CI red, an unresolved merge of main, a merged head the gate failed on) and rebase (a code conflict, the gate or pre-push check failing on the merged head), each by spawning the named role, then land, render and rhiza session done; plus handback_cap and that the steps after the merge never hand back.

## Links

### Out

- decided-by [conduct-guide-three-topics](../../../../decision/conduct-guide-three-topics.md)
- decided-by [conduct-handback-routes](../../../../decision/conduct-handback-routes.md)
- decided-by [conduct-handback-routes-ship-rebase](../../../../decision/conduct-handback-routes-ship-rebase.md)

### In

- verifies [conduct-handback-topic-verified](../../../../evidence/conduct-handback-topic-verified.md)
- verifies [conduct-handback-topic-qa-1](../../../../evidence/conduct-handback-topic-qa-1.md)
- verifies [conduct-handback-topic-qa-2](../../../../evidence/conduct-handback-topic-qa-2.md)

## History

- `01M7MJN83Y6EF727C2T0AF3Y5D` The rhiza-conduct skill and its guide topics are enough on their own: a conductor outside this repo has no Rhiza source to read, so everything a conductor looked up in the code here (step roles, handback routes, workflow-new defaults, check roles) is in the guide
- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict
