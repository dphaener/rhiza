# resume-only-in-flight-plan-interview

source · `01M62RGWE4EKV11H26KTW63G8Y`

## quotes

```json
[
  {
    "by": "agent (accepted by the owner)",
    "text": "Show in inbox (Recommended)",
    "where": "plan interview, skipped-run question"
  },
  {
    "by": "agent (accepted by the owner)",
    "text": "Refuse too (Recommended)",
    "where": "plan interview, manual-resume question"
  }
]
```

## raw_sha256

5434da4079326cc62a662e83e3fc3c4b1b071d6e35ff0865be1934d89c78dfb9

## session

16eb7927-abe1-4ab2-84b7-01a185856b86

## summary

The daemon resumes an interrupted run only while its change set is accepted. The owner picked, from the agent's options: an interrupted run on a change set that is no longer accepted (landed, rejected or withdrawn) shows in the inbox with the change set's state and the next step (rhiza run abandon); and rhiza run resume refuses such a run too, rather than only the daemon. No typed free answers; both decisions are picks of the agent's options, accepted by the owner.

## title

Resume only in-flight runs, planned with the owner

## Links

### In

- derived-from [skipped-run-in-inbox](../decision/skipped-run-in-inbox.md)
- derived-from [manual-resume-refuses](../decision/manual-resume-refuses.md)

## History

- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
