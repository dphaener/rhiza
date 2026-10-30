# usage-limit-resume-same-session

decision · `01M650608BEF1XS31K5BFKJHPM`

## reasoning

The question says the transcript persists on disk, and the Claude adapter already starts each session with a stored --session-id (claudecode.go:57). Resuming that session keeps the agent's context and the work it already committed, so the run does not redo finished steps. A fresh session would lose that context and need a re-briefing. The fresh-session fallback covers the case where the transcript is gone. Nothing in the adapter needs a new ID scheme for this. It departs from engine-orchestrates (agent steps are short, fresh sessions) only for a session that hit the usage limit, so the rule stays for every other step.

## statement

Resume a usage-limited run by continuing the same agent session with the Claude CLI's resume flag, keyed on the session ID the run already stores; start a fresh session only when the transcript is missing. This is the one exception to engine-orchestrates: a step's agent session stays open across a usage-limit pause

## wording

(assumed, would have escalated)

## Links

### Out

- relates-to [engine-orchestrates](engine-orchestrates.md) — Exception to engine-orchestrates: a usage-limited step resumes its agent session after the pause instead of starting a fresh one
- answers [question-6e2b93b5-1](../question/question-6e2b93b5-1.md)
- derived-from [usage-limit-pause-interview](../source/usage-limit-pause-interview.md)

### In

- decided-by [usage-limit-resume](../capability/rhiza/workflows/workflow-runs/run-resume/usage-limit-resume.md)
- decided-by [usage-limit-fresh-session](../capability/rhiza/workflows/workflow-runs/run-resume/usage-limit-fresh-session.md)

## History

- `01M61R3RG70WZD5EDH293ZGVWN` Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.
