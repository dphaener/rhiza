# m40-interventions-plan-research

source · `01M7R2ESXAWZBHM2HKWK57EHMM`

## quotes

```json
[
  {
    "by": "owner",
    "text": "the overall idea here is that the daemon should be able to execute any workflow from start to finish without intervention. sure there may end up being errors from time to time that require intervention but for the most part it should just work. intervention should be the exception.",
    "where": "hands-off milestone interview, 2026-11-18"
  }
]
```

## raw_sha256

869ba6c2730cd3f921a50907f41b7dc8aa5e668f126e32f70861ad3586f1bc5a

## summary

Autonomous plan of change set 01M7NWFV: no owner interview. Material: the milestone-40 hands-off interview notes (the orchestrator's record of the owner's goal and the stopgaps it ran by hand) and a read of the code. Findings: session.message (internal/api/session_message.go) journals nothing; a plain run.resume journals nothing and a manual --retry journals the same TypeRetried as a watchdog retry; run.abandon journals a finished entry; run.approve only acts at a waiting signoff step (the workflow asking the owner); an owner-rule accept/reject journals nothing in a run; nothing detects commits on a run's branch that the run did not make (run commits carry Leaf:, Review:, Change:/Actor: trailers); rhiza run status builds RunInfo (internal/api/run.go) and rhiza daemon status builds DaemonStatusResponse with DaemonRun rows (internal/api/daemon.go). The orchestrator's later evidence names a hand rhiza run start that retried a failed run's close-out as an intervention. The milestone's measure is consecutive daemon runs finishing with zero interventions.

## title

Interventions plan research (autonomous conduct run 01M7R281277T47T6RJ7MF8QJXD)

## Links

### In

- derived-from [intervention-one-entry-type](../decision/intervention-one-entry-type.md)
- derived-from [intervention-is-a-person-command](../decision/intervention-is-a-person-command.md)
- derived-from [hand-commit-by-trailer](../decision/hand-commit-by-trailer.md)
- derived-from [daemon-counts-clean-streak](../decision/daemon-counts-clean-streak.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted
