# plan-resume-landed-record-interview

source · `01M6GD2FJ12DMGE6Y345CJDJ53`

## quotes

```json
[]
```

## raw_sha256

8bb6faec213a21030087408af9af42b5dfc17dd4d0e5e736aa6f2b11cb4bb766

## session

b693317e-dc4d-40eb-b2ec-01a1a0cbe9f7

## summary

Autonomous plan run of change set 01M6GCQA. The planner read internal/api/run.go: runRunResume refuses any unfinished run whose change set is in a terminal state (changeset.State.Terminal(): landed, rejected, withdrawn, superseded). A plan run lands its record during the interview, so a plan run that the gate stops on open questions (capability plan-run-stops-open) can never resume. The planner raised two questions (which intent to serve; how narrow the exemption is) and the arbiter answered with Decision plan-resume-landed-record-milestone-6: serve milestone-6; let through only plan-workflow runs whose change set is landed, keeping the refusal for rejected, withdrawn and superseded change sets and for every other workflow. One leaf was planned under plan-gate, beside plan-run-stops-open.

## title

Plan of the Fix: resume a plan run whose record landed (autonomous)

## Links

### In

- derived-from [plan-resume-landed-record-milestone-6](../decision/plan-resume-landed-record-milestone-6.md)

## History

- `01M6GCQACNNTMW5B57HW8N7K50` Fix: rhiza run resume refuses a plan run whose record change set landed (internal/api/run.go ~792, the Terminal() guard). Every plan run lands its record during the interview, so a plan run stopped at the gate on open questions can never resume. Exempt plan-workflow runs from that guard
