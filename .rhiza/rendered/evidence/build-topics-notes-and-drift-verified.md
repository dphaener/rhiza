# build-topics-notes-and-drift-verified

evidence · `01M7KYBS5H07KW28K8XDKY0Y2Z` · actor [claude-code-1b0a4731](../actor/claude-code-1b0a4731.md)

## result

PASS

## statement

At 5529d114f: go test -count=1 -trimpath ./internal/guide/... ok (TestBuilderText note+drift phrases, TestVerifierText note phrases plus the negative check that 'edit code' appears only as 'Never edit code' and no rhiza set/fix commit, TestRunNoteOnlyInBuildTopics over check-plan/code/qa/compliance/quality/model, critic and plan-critic, TestBuildRetryText 'quoted verbatim in the prompt', TestConduct retry-brief phrases, topic token caps); go test -trimpath ./internal/api -run TestGuideInstall ok; make self-check exit 0 and rhiza ship generate / rhiza render with no diff in a clean detached worktree (in the run worktree self-check only flagged the verifier's own untracked actor record). Diff touches only internal/guide text and tests, plus node.update entries for builder-topic, verifier-topic, build-retry-topic and conduct-skill matching the new text; builder text names no Go/make command; no stubs or TODOs.

## Links

### Out

- verifies [build-topics-notes-and-drift](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-guide/build-topics-notes-and-drift.md)

## History

- `01M76XT9E2N9ZYMMB8DMTMKA3Q` A leaf's retry prompt carries the verifier's findings verbatim and says to inspect the working tree and address every finding as written; builders and verifiers append one-line learnings to the run's notes, which later leaves' prompts carry; builders update landed nodes their leaf changes or raise a question. Evidence: 113 of 368 leaves needed a second attempt, 36 model-check majors were drift (prompt-gap recs 9-11).
