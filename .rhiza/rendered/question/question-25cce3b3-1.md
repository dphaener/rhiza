# question-25cce3b3-1

question · `01M6JZH5R5E11NJ63V7Z66DSRH`

## session

25cce3b3-9133-4ff2-a0fa-01a1a5f70158

## statement

Change set 01M6JZD: an interrupted plan run whose record change set has landed (a milestone plan lands its record mid-interview) can now be continued by `rhiza run resume`, but the daemon's Resumable (internal/daemon/pick.go) resumes only runs whose change set is accepted, and the inbox (internal/api/inbox.go ~249) lists the run as skipped with `rhiza run abandon`. What should the daemon do with such a run? / Change set 01M6JZD: the inbox also hides a failed or gave-up run when its change set is terminal (gone() in inbox.go), which hides a failed plan run on a landed record even though `rhiza run resume --retry` works on it. Should this change set also stop hiding those?

## Links

### In

- answers [decision-plan-landed-resumable](../decision/decision-plan-landed-resumable.md)

## History

- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
