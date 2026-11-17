# rejoined-conductor-done-mark-plan

source · `01M7MD6PJZ98G485C02QK5Q0J0`

## quotes

```json
[]
```

## raw_sha256

d054e54b172419e3a3fdab4269e67876b25efdee51ced6f34ae5a35e196a12db

## session

09f76cee-7955-4259-9b5e-01a1e8d192bb

## summary

Autonomous plan run of a bug request; no owner, no questions raised. Research: stepRunner.Handback (internal/api/run.go) resumes an ended session through resumeSession, which clears the record's whole Status, but messages a live session through messageSession and leaves Status.Done (the first visit's rhiza session done mark) in place. AgentStatus then reports SaidDone for the rejoined visit, so a step waiting for done (ship-change's conduct, until_when = done) can pass, and anything closing done sessions can close it, before the session has worked the hand-back. The plan clears the done mark when a hand-back messages a live session, before the message is typed, and leaves the resume path, the asked mark and the other status reports as they are.

## title

Plan of the rejoined conductor done-mark fix (autonomous plan run)

## Links

### In

- derived-from [handback-clears-done-mark](../decision/handback-clears-done-mark.md)

## History

- `01M7MD2K2BE981QPPF5HE8M551` Fix A rejoined conductor session keeps its first visit's done mark, so after a gate hand-back the session reads as done before it has worked the hand-back (anything closing done sessions closes it mid-visit)
