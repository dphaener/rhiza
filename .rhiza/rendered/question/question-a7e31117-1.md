# question-a7e31117-1

question · `01M6MYM641F6JEE2SKH35DXQGH`

## session

a7e31117-48f2-4f83-8965-01a1a9e80a6a

## statement

When does the daemon start a queued milestone request's plan run? Today every request waits for a free run slot (max runs), oldest first, so a milestone interview can sit behind build runs. / How is the owner told a milestone interview waits on them? Today a notification goes out only when the agent's first AskUserQuestion blocks, which can be minutes after the tab opens, and the inbox lists nothing before that. / Does a usage-reserve hold (the governor holding new launches) also hold a milestone interview's start?

## Links

### In

- answers [milestone-interview-starts-at-once](../decision/milestone-interview-starts-at-once.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
