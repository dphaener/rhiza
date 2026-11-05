# milestone-planning-always-plan

source · `01M6MYEZJY2J3AEG84YZ9YKBR3`

## quotes

```json
[
  {
    "by": "owner",
    "text": "milestone planning should never be autonomous even in autonomous mode. that always requires a human.",
    "where": "the change set's intent, 2026-11-05"
  }
]
```

## raw_sha256

dcb9e9c0f236b97e22b5c225dbeceb13facedc68ac59a09a5c5e197a54e2a9b6

## session

58fb9cbb-b516-4d6f-a418-01a1a9e44747

## summary

Autonomous plan run of change set 01M6MY8B. The owner's request: milestone planning never runs autonomously, even with plan autonomy auto. The planner found autonomy enters plan runs in two places: the daemon's daemonStartPlan (internal/api/request.go), which passes --autonomous whenever rhiza.autonomy.plan is auto, and run start --autonomous (internal/api/run.go). It asked which intent the change set serves and how to enforce the rule. The arbiter decided (milestone-plan-never-autonomous): serve milestone-22; the daemon never passes --autonomous for a milestone request; run start refuses --autonomous on a plan run that plans a milestone (--milestone, a milestone request, or an existing record change set), and so does resume; feature and bug requests keep following rhiza.autonomy.plan. Two leaves follow, nothing left open.

## title

Milestone planning always needs a human, planned autonomously

## Links

### In

- derived-from [milestone-request-owner-plan](../capability/rhiza/orchestration/daemon/plan-requests/request-queue/milestone-request-owner-plan.md)
- derived-from [milestone-plan-refuses-autonomous](../capability/rhiza/workflows/autonomous-runs/milestone-plan-refuses-autonomous.md)

## History

- `01M6MY8BMPFK9ZWD3TDQHKT64J` Milestone planning always needs a human: a plan run for a milestone (a record change set, or a request of kind milestone) never runs autonomously, even when the daemon or run is autonomous (rhiza.autonomy.plan=auto); it waits for the owner's interview instead of the arbiter answering. Owner, 2026-11-05: "milestone planning should never be autonomous even in autonomous mode. that always requires a human."
