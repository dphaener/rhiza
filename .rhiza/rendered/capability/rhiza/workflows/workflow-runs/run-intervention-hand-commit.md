# run-intervention-hand-commit

capability · `01M7R2FY8VC5AK2F7BWWFFHW55` · parent [workflow-runs](../workflow-runs.md)

## acceptance

go test ./internal/api -run HandCommit (or ./internal/run): a branch with a Leaf:-trailered commit, a Review: <run> commit, an Actor: commit of a run session and a merge of main shows 0 hand commits; a plain commit made by hand shows one hand-commit intervention naming it, and a second scan adds none; (failure case) a commit whose Review: trailer names another run counts as a hand commit

## files

internal/api/run_handcommit.go, internal/api/run_handcommit_test.go, internal/run/engine.go

## must_not

Block, fail or change the run because of a hand commit; scan outside the run's own branch

## statement

The engine, on entering each step and before the run finishes, journals one hand-commit intervention (with the commit) for each commit on the run's branch since it started that the run did not make, never twice for one commit

## Links

### Out

- depends-on [run-intervention-journal](run-intervention-journal.md)
- decided-by [hand-commit-by-trailer](../../../../decision/hand-commit-by-trailer.md)

### In

- verifies [run-intervention-hand-commit-verified](../../../../evidence/run-intervention-hand-commit-verified.md)
- verifies [run-intervention-hand-commit-qa-1](../../../../evidence/run-intervention-hand-commit-qa-1.md)
- verifies [run-intervention-hand-commit-qa-2](../../../../evidence/run-intervention-hand-commit-qa-2.md)

## History

- `01M7NWFV54PDK3XDSGH07HK0K7` The run journal records every outside intervention in a run (a session message from outside the run, a manual resume, retry or abandon, an approve of anything other than a step the workflow itself asks the owner to approve, a hand commit on the run's branch), and rhiza run status and rhiza daemon status count them, so runs finished with no intervention can be counted; approving a workflow's own signoff (or other owner) step is the workflow working as written, not an intervention, and is not counted
