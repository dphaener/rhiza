# inbox-rest-cut-verified

evidence · `01M7PP0VTRDEYQCMBAGC9STYP2` · actor [claude-code-2238313f](../actor/claude-code-2238313f.md)

## result

PASS

## statement

Verified at 3bf00d314 in a detached worktree. go test -trimpath -count=1 -p 2 ./internal/api -run 'Inbox': ok 13.5s; the Checks: Proposal|Uncovered|NeedsYou|Daemon.*Accept|QuestionAdvice|Advice run: ok 6.6s; ./internal/index: ok 40.3s (new TestQueryChangeSetStates); ./cmd/rhiza -run TestScripts/inbox: ok (20 scripts); internal/model untouched. RHIZA_CPUPROFILE, 6 warm rhiza inbox runs merged per binary in the main checkout, parent d97acc451 -> commit: run.Journal.State 19.01% -> 5.38%, api.uncoveredProposals 7.44% -> gone from the listing, inboxFold.load 4.55% -> 2.69% (index.ChangeSetStates 2.69%), Index.ChangeSets 7.85% -> 6.45%, questionAdvice 2.89% -> gone; total samples 2.42 -> 1.86 s. /usr/bin/time rhiza inbox in the main checkout, one warm call each, base 5d2596d and commit alternately x5, load ~3.5: base user+sys 1.09 1.10 1.02 1.09 1.07 (median 1.09 s), commit 0.34 0.31 0.32 0.33 0.32 (median 0.32 s): 29.4%, target at most 31%, met. Failure case: parent and commit inbox --json in the main checkout cmp 0 twice; in a scratch worktree after rhiza note add on an accept-listed change set, cmp 0; after rhiza withdraw of it, both binaries drop it (10 items), cmp 0. Fast path in inboxFold.load matches changeSetInfos with no pending entries (same stored rows). rhiza ship generate exit 0, rhiza render no diff, make self-check exit 0 in the clean checkout. No TODOs; no format or threshold changes. Note: internal/api/proposals.go is changed outside the leaf's files list (inbox code the profile names, said so in the commit); accepted as within the statement's internal/api scope.

## Links

### Out

- verifies [inbox-rest-cut](../capability/rhiza/orchestration/inbox/inbox-rest-cut.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
