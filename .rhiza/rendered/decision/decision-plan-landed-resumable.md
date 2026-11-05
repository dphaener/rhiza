# decision-plan-landed-resumable

decision · `01M6JZJYNR0W6K12FR5VKYQW0X`

## reasoning

The change set's intent says the inbox and the daemon should treat this run as resumable. run.go's resume check already makes this exact exception: a milestone plan lands its record during the interview, so for a plan run landed is a normal state partway through, not an ending. The daemon and the inbox are the last places that still use the old rule (terminal means gone), and one rule used in all three places keeps them consistent. Before this, the inbox told the owner to abandon a run that `rhiza run resume` happily continues, and it hid a failed run that `--retry` can fix. Both are wrong next steps. The owner's goal is for the daemon to do the orchestration, so it should give the run the same one automatic resume as any other interrupted run; the existing gave-up item already covers a second interruption. Keeping the exemption to plan runs on landed change sets, and handling both parts of the question in this change set, keeps the fix small and stops it from reviving runs whose change set really did end.

## statement

A plan-workflow run whose change set has landed counts as still in flight, everywhere `rhiza run resume` already lets it through (internal/api/run.go). The daemon's Resumable (internal/daemon/pick.go) gives it the same one resume as an interrupted run on an accepted change set. The inbox stops listing it as skipped with `rhiza run abandon`. If that one resume is interrupted again, it shows as gave-up with `rhiza run resume <id>`. Yes, this change set also stops hiding these runs: gone() in internal/api/inbox.go does not count a plan run's landed change set as terminal, so a failed plan run on a landed record shows as failed with `rhiza run resume --retry <id>`. The exemption covers only the plan workflow and only the landed state. A plan run whose change set was rejected, withdrawn or superseded, and every non-plan run on a terminal change set, are skipped or hidden as before.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-25cce3b3-1](../question/question-25cce3b3-1.md)
- derived-from [plan-landed-resumable-interview](../source/plan-landed-resumable-interview.md)

### In

- decided-by [daemon-resumes-landed-plan-run](../capability/rhiza/orchestration/daemon/daemon-resumes-landed-plan-run.md)
- decided-by [inbox-landed-plan-run-in-flight](../capability/rhiza/orchestration/inbox/inbox-landed-plan-run-in-flight.md)
- decided-by [daemon-resumes-accepted-only](../capability/rhiza/orchestration/daemon/daemon-resumes-accepted-only.md)
- decided-by [daemon-resume-once](../capability/rhiza/orchestration/daemon/daemon-resume-once.md)
- decided-by [unaccepted-run-in-inbox](../capability/rhiza/orchestration/inbox/unaccepted-run-in-inbox.md)
- decided-by [inbox-clears-stale](../capability/rhiza/orchestration/inbox/inbox-clears-stale.md)

## History

- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
