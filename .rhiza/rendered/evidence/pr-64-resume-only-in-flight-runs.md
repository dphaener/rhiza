# pr-64-resume-only-in-flight-runs

evidence · `01M6337AGWXXAE5EZ0Z8D15JS8` · actor [darin-haener](../actor/darin-haener.md)

## head

f4bc9ad4e8acc0e53b139efae98b35efd1c92818

## pr

https://github.com/dphaener/rhiza/pull/64

## result

1 check passed: check SUCCESS

## statement

Pull request #64 (https://github.com/dphaener/rhiza/pull/64) passed CI at f4bc9ad

## status

verified

## Links

### Out

- verifies [daemon-resumes-accepted-only](../capability/rhiza/orchestration/daemon/daemon-resumes-accepted-only.md)
- verifies [unaccepted-run-in-inbox](../capability/rhiza/orchestration/inbox/unaccepted-run-in-inbox.md)
- verifies [manual-resume-refuses-unaccepted](../capability/rhiza/workflows/workflow-runs/manual-resume-refuses-unaccepted.md)
- verifies [daemon-resume-once](../capability/rhiza/orchestration/daemon/daemon-resume-once.md)
- verifies [inbox-list](../capability/rhiza/orchestration/inbox/inbox-list.md)
- verifies [inbox-notify](../capability/rhiza/orchestration/inbox/inbox-notify.md)
- verifies [run-resume](../capability/rhiza/workflows/workflow-runs/run-resume.md)

## History

- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
