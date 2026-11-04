# actor-once-per-run-plan

source · `01M6M2Z6M2VESQQ0T25NT0YTYF`

## quotes

```json
[]
```

## raw_sha256

30465c275482099475f75a199b8708678ed01a0552277aa3f11875b2983db9c7

## session

7346bd67-434f-4dba-9398-01a1a82c3dce

## summary

Autonomous plan run, no owner present. The planner traced both duplicate actor records (claude-code-a463ece8 and claude-code-83858632) to the same cause: session start (recordAgentActor) writes the agent's actor record only in the run worktree, uncommitted, while a verifier works in its own checkout that lacks it, so the agent hand-added 'add actor <handle>' there and both records were later committed. The planner raised the intent and the fix as one question; the arbiter decided (actor-record-once-per-run): the change set serves m4-runs-hold-up; session start copies the same record change set's log entries into a session's own checkout; recordAgentActor reuses an existing actor whose session field matches; rhiza check and land do not learn to merge duplicate actor records. Two leaves under actor-stamping were planned from that.

## title

Plan: Fix an agent session's actor record (autonomous)

## Links

### In

- derived-from [actor-record-once-per-run](../decision/actor-record-once-per-run.md)

## History

- `01M6M2R7RJCMNMMKH0KS2364BK` Fix An agent session's actor record is written once per checkout, not once per run, so a ship run's branch ends up with two record change sets adding the same actor handle and its review/code step fails rhiza check with 'duplicate handle'. Seen twice on 2026-11-04: usage governor (claude-code-a463ece8) and text references (claude-code-83858632: recorded uncommitted in the run worktree as change set 01M6KE23, then again in the verifier checkout as 01M6KFMJ, both committed). Make actor stamping idempotent across a run's worktree and checkouts (look up the handle in the branch, the run worktree and main before adding), or have rhiza check/land merge identical actor records
