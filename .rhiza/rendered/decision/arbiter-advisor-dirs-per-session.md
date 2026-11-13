# arbiter-advisor-dirs-per-session

decision · `01M7AMH28YP1KWX1S2Y6C1JGC8`

## reasoning

The finding is real. In internal/run/engine.go, advise() calls StartAdvisor with the same asker session each time that session blocks (line 1090), and undecidedAsks lets a retry have the arbiter decide the same question anew. So advisor-<asking session id> and arbiter-<question id> can both collide, and the second session would overwrite the first's brief.md and slice.json. The arbiter is Definition.Arbiter and the advisor is advisorStep. Neither is entered as a visit (the entered entries carry Visit only for workflow steps), so <step>-<visit> has no visit number to give them. A counter would need the journal, which StartArbiter and StartAdvisor in internal/api/run.go do not see. The session's own ID is unique by construction and available in time: runSessionStart in internal/api/session.go allocates sid via newSessionID before it launches. It is also deterministic under RHIZA_TEST_SEED for tests. This keeps the decision's goal (no two sessions share a directory, durable beside the journal, outside any worktree) and fixes both names in one rule, without changing either prompt.

## statement

The arbiter and advisor sessions are not workflow step visits (no entered entry, no visit number), so they follow brief-files-and-invoke-placement's one-directory-per-session rule with their own key instead of <step>-<visit>: each gets <commondir>/rhiza/runs/<run id>/arbiter/<session id>/ or <commondir>/rhiza/runs/<run id>/advisor/<session id>/, keyed on the session's own ID, which runSessionStart allocates before the agent starts. The directory is filled with brief.md and slice.json after that ID is allocated and before the launch. Two advisors for one asking session, or an arbiter re-run on the same question after a retry, so never share a directory. arbiter-advisor-brief's directory names change to these, and so does its acceptance: a second advisor started for the same asker gets a separate directory, and the first one's files are left alone.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7ah6rhkmzrahgjkjfm4vtnt-r3-f2](../question/review-01m7ah6rhkmzrahgjkjfm4vtnt-r3-f2.md)

## History

- `01M79KGHYF3272SZ9H9Y12FWF4` Every agent session a run starts gets one gate-checked contract and its context in files: brief.md, slice.json and, on a hand-back, handback.md, with templated invoke and handback text on agent steps
