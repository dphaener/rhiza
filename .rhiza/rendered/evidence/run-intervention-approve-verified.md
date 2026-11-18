# run-intervention-approve-verified

evidence · `01M7R4V3EPW4M0NJYWHNZRZ5XA` · actor [claude-code-e68f0837](../actor/claude-code-e68f0837.md)

## result

PASS

## statement

Verified bebbb7d91 under the capped scope with RHIZA_/SHIP_/REX_ unset. Exit 0: go test ./cmd/rhiza -run TestScripts/run-intervention-approve$; the 9 named neighbour scripts (autonomy-accept-rule, owner-run-gate-accept, inbox-owner-run-accept, run-worktree-carries-accept, run-status-interventions, run-signoff-workflow-owner, review-signoff, review-signoff-autonomous, lifecycle); go test ./internal/api -run 'Accept|Reject|Lifecycle|Gate|Intervention|Signoff|Approve'; go vet ./internal/api; rhiza ship generate (0 files staged); rhiza render. The Checks line for the 64 accept/reject scripts was a placeholder, so I rebuilt the list as the scripts containing 'rhiza accept|reject' and 'run start' (exactly 64). They passed twice; one earlier run printed FAIL with no output kept, under load average 8.7 on 4 cores, and did not happen again. The script covers each acceptance point: an owner accept under auto accept makes intervention_count 1 with kind approve, step hold and actor test; an owner reject of a live run journals approve too; a dry run journals none; a run whose workflow declares accept owner journals none; an accept with no live run leaves every run's count unchanged; rhiza run approve at a waiting signoff gives succeeded with intervention_count 0. From the diff: journalApprove runs only when rule == owner and not on a dry run. The workflow's own gate (run_review.go:1843) always accepts under rule autonomous or autonomy, so it never journals. rhiza run approve (run_signoff.go) is untouched and never calls journalApprove. A journal error only prints a warning, so accept, reject and approve behave as before. Nothing changed outside the leaf's files (the test workflow's own autonomy setting covers run_signoff.go). No TODOs or stubs.

## Links

### Out

- verifies [run-intervention-approve](../capability/rhiza/workflows/workflow-runs/run-intervention-approve.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted
