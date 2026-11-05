# milestone-request-launch-plan

source · `01M6MYRT7GAXXCPJZZAWQN5FAW`

## quotes

```json
[]
```

## raw_sha256

be4603c3027fde058d709ae7304d3912e61822c0a078d15d6d8cc0eb9a9d0f97

## session

a7e31117-48f2-4f83-8965-01a1a9e80a6a

## summary

Autonomous plan run of change set 01M6MYFH4QYGSKDATBH509SEJF (A milestone request launches an interactive planning session). The planner read internal/api/request.go (daemonPlanRequests starts requests only while running < max, oldest first, held by the usage governor), internal/api/inbox.go (a non-autonomous run's question item appears only once an AskUserQuestion blocks) and internal/api/daemon.go (daemonNotify). The daemon already starts a milestone request's plan run and its interview opens a Rex tab; what is missing is starting it at once and telling the owner at launch. Three questions (start timing, notification, usage hold) went to the arbiter, which answered with Decision milestone-interview-starts-at-once: start at the next poll ahead of every request and past run slots, one milestone interview live at a time; waiting status at launch plus an inbox item from launch until the run ends; not held by a usage-reserve hold, still held by pause and disk-low. The never-autonomous part is the sibling change set 01M6MY8BMPFK9ZWD3TDQHKT64J. The change set serves milestone-25. No owner was present, so no owner words are quoted.

## title

Planning the milestone-request launch, autonomous run

## Links

### In

- derived-from [milestone-interview-starts-at-once](../decision/milestone-interview-starts-at-once.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
