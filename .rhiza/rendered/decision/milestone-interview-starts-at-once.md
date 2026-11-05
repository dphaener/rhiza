# milestone-interview-starts-at-once

decision · `01M6MYP9D5XEFXQ5DHFXZQ8JDZ`

## reasoning

The change set's intent is that a milestone request 'automatically launches' an interactive session and 'notifies the owner that it is waiting', and milestone planning always needs a human: an interview queued behind build runs (daemonPlanRequests only starts requests while running < max, after ship picks and queued starts) defeats that, and the session spends most of its time waiting on the owner, so it costs little compute and need not compete for a slot; one at a time because there is one owner to interview. Notifying only at the first AskUserQuestion leaves the owner unaware of a tab that already waits on them, and Rex has no notify API (notifications follow the program's OSC 7501 status), so the status must say waiting at launch and the inbox must list it. internal/api/usage_hold.go says the governor holds what the daemon would launch while 'live runs and the owner's own starts are untouched'; a milestone request is the owner's own, attended start, and the reserve exists to keep usage for exactly that kind of owner-driven work, so it is not held. Pause is the owner's explicit stop and disk-low is a safety guard (a run needs a worktree), so neither is changed by this answer.

## statement

A queued milestone request's plan run starts at the daemon's next poll, ahead of every other queued request and without waiting for a free run slot: it is the owner's own interactive start, not a pick. It still counts as a running run once started, so later picks see it, and only one milestone interview is live at a time; a second milestone request waits for the first to end. The owner is told the moment the interview's tab opens: the plan session's status goes to waiting-on-the-owner at launch (so Rex's notification fires from the OSC 7501 status) and the inbox lists the live milestone interview as an item from launch until the run ends, not from the first AskUserQuestion. A usage-reserve hold does not hold a milestone interview's start; it keeps holding build, review and other autonomous launches. Pause and the disk-low guard are unchanged and still hold it.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-a7e31117-1](../question/question-a7e31117-1.md)
- derived-from [milestone-request-launch-plan](../source/milestone-request-launch-plan.md)

### In

- decided-by [milestone-request-starts-at-once](../capability/rhiza/orchestration/daemon/plan-requests/milestone-request-starts-at-once.md)
- decided-by [milestone-request-past-usage-hold](../capability/rhiza/orchestration/usage-governor/milestone-request-past-usage-hold.md)
- decided-by [milestone-interview-waiting-at-launch](../capability/rhiza/orchestration/notifications/milestone-interview-waiting-at-launch.md)
- decided-by [inbox-milestone-interview](../capability/rhiza/orchestration/inbox/inbox-milestone-interview.md)
- decided-by [request-queue](../capability/rhiza/orchestration/daemon/plan-requests/request-queue.md)
- decided-by [request-slot-kept](../capability/rhiza/orchestration/daemon/plan-requests/request-slot-kept.md)
- decided-by [usage-hold](../capability/rhiza/orchestration/usage-governor/usage-hold.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
