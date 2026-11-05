# plan-landed-resumable-interview

source · `01M6JZM374G612880NXAE1GK0T`

## quotes

```json
[]
```

## raw_sha256

a97649773abf2e1cd7923442c694d50c6ff15cabe6d06b7d6a0d3f3e7d316676

## session

25cce3b3-9133-4ff2-a0fa-01a1a5f70158

## summary

Autonomous plan run, no owner. The planner found that rhiza run resume (internal/api/run.go) already continues an interrupted plan run whose record landed, while the daemon's Resumable (internal/daemon/pick.go) resumes only runs on accepted change sets and the inbox (internal/api/inbox.go) lists such a run as skipped with rhiza run abandon, and its gone() hides failed and gave-up plan runs on a landed record. Two questions went to the arbiter: should the daemon resume such a run, and should gone() stop hiding them. The arbiter ruled (decision-plan-landed-resumable, assumed, would have escalated): a plan-workflow run on a landed change set counts as in flight; the daemon gives it its one resume; the inbox stops listing it as skipped, shows gave-up with rhiza run resume after a second interruption, and shows a failed one with rhiza run resume --retry. Only the plan workflow and only landed are exempt. The plan also replaces leafless draft 01M6H5VQ231Z8V (linked relates-to, withdrawn) and serves milestone-25.

## title

Plan of 01M6JZD: plan runs on a landed record are resumable (autonomous)

## Links

### In

- derived-from [decision-plan-landed-resumable](../decision/decision-plan-landed-resumable.md)

## History

- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
