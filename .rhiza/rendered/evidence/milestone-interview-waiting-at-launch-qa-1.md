# milestone-interview-waiting-at-launch-qa-1

evidence · `01M6N3HX43WFBZ0VBC3N7VJXKS` · actor [claude-code-23c10d47](../actor/claude-code-23c10d47.md)

## result

PASS

## statement

QA at dea4e5b6 with a scratch build (go build -trimpath) in a scratch git repo, fake rex that runs the tab command and a fake claude that sleeps: rhiza run start --foreground --workflow plan --new 'Backlog empty' --milestone gave OSC 7501 states working, blocked, done on stderr, the blocked msg decoding to '<cs>: interview waits on you in tab <session>: the milestone interview waits on you', emitted before the agent asked anything; a feature plan run and a --milestone --autonomous run gave only working, done. The run's journal carries milestone_plan:true only for the milestone run.

## Links

### Out

- verifies [milestone-interview-waiting-at-launch](../capability/rhiza/orchestration/notifications/milestone-interview-waiting-at-launch.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
