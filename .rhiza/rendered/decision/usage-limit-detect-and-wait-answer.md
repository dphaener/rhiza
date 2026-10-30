# usage-limit-detect-and-wait-answer

decision · `01M65KAXW8HAESW7N737Q9YHPH`

## options

```json
[
  "reset time from the limit message, else a fixed 15 minutes",
  "a fixed interval only"
]
```

## statement

The Claude adapter recognises the usage limit when the last assistant line of the session's transcript is Claude Code's synthetic API-error message (isApiErrorMessage, model '<synthetic>') whose text says 'usage limit reached' (the reset as a Unix time after a bar) or 'hit your limit'; the phrases in the model's own turns, its prose or a tool call's input, do not count. While the agent runs, the run reads the transcript every 30 seconds and stops an agent left open at its limit; after the agent ends, a limited session waits until the reset time the message gives, else a fixed 15 minutes, then resumes. While it waits the session is recorded as working with detail 'paused at its usage limit until <reset>', not idle: a Stop hook reporting idle for any session whose transcript shows the limit records it as working, 'paused at its usage limit', so a run and rhiza session show see it paused rather than ready at its prompt. One session pauses at most 24 times: still limited after that, the session fails with that reason instead of passing.

## wording

the implementing agent's choice, made by the review-cycle-2 fixer to answer the confirmed qa and model findings and extended by the review-cycle-3 fixer (synthetic-message detection, the paused state) for the confirmed quality, qa and model findings; not the owner's words

## Links

### Out

- answers [usage-limit-detect-and-wait](../question/usage-limit-detect-and-wait.md)

### In

- decided-by [usage-limit-resume](../capability/rhiza/workflows/workflow-runs/run-resume/usage-limit-resume.md)

## History

- `01M61R3RG70WZD5EDH293ZGVWN` Usage-limit pause: a run whose agent session hits the Claude usage limit waits for the limit to reset and then continues the session, instead of failing.
