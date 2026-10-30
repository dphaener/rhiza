# roadmap-app-accept-answer-verified

evidence · `01M675SJEZKVA6JJNH68GNDSZW` · actor [claude-code-b5fd5234](../actor/claude-code-b5fd5234.md)

## result

PASS

## statement

At 270574d: go test ./cmd/rhiza -run TestRoadmapApp -count=1 -v passes both tests. TestRoadmapAppAcceptsAProposedChangeSet builds a temp repo with a proposed change set, presses a, asserts the view shows the What-if (proposed -> accepted) and y/n with the state still proposed, presses y, asserts change show reports accepted and the list reloads. TestRoadmapAppAnswersABlockingQuestion opens the detail of a change set blocked by a question, checks rhiza ready shows it blocked, presses enter, types 'Use CSV', enter, then asserts decision answer-which-format quotes the text and links answers, rhiza ready reports the change set ready, and rhiza render + rhiza check succeed. Diff read: internal/tui/roadmapapp/answer.go writes kind=decision with quote=the typed words (owner) and an answers link in that change set. Attempt 1's finding (stale rendered actor file) is fixed: in a clean detached worktree of 270574d, make fmt-check, go test ./... -count=1 and make self-check all pass. In this run checkout, self-check flags only this verifier session's own uncommitted actor record (claude-code-b5fd5234), which isn't part of the commit.

## Links

### Out

- verifies [roadmap-app-accept-answer](../capability/rhiza/roadmap-app-accept-answer.md)

## History

- `01M673223RSVC166NJVVA13X75` The roadmap app is a cursor list of milestones and change sets that the owner acts on: open a change set's detail, plan it, ship it, accept it, answer its blocking questions, or jump to its run's Rex tab
