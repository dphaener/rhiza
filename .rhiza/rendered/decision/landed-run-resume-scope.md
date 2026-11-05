# landed-run-resume-scope

decision · `01M6J9GAVT8PG9PXD2VN09EEQV`

## reasoning

m4-runs-hold-up says 'A build, review or ship-change run finishes without hand repair', and a stranded ship-change run that has to be abandoned and closed out by hand is exactly that kind of hand repair. It is the only landed intent the fix serves. The question says 01M6H5VQ already covers the plan-run half (inbox and daemon), and internal/api/run.go already exempts planWorkflow from the landed refusal in run resume. Covering plan runs here too would put two drafts on the same code. For what counts as resumable: ship-change.toml runs land, then wait-landed, landed-green, rebase-merge, merge, wait-main, main-green, close-out and ship. A change set can only have landed through a run once that run reached land. rhiza ship land reports landed=false and does nothing when the change set already landed (ShipLandResponse.Landed), so re-running land is safe, and the steps after it are idempotent waits, checks or merges. 'The run reached land' is a test anyone can read off the journal. It also keeps run-resume-landed-plan.txtar valid, since that script expects a ship-change run whose change set landed but whose journal never reached land (BUILDLAND01) to be refused. The daemon side is the same test, swapped in for the bare accepted(st.Change) check in internal/daemon/pick.go Resumable

## statement

Change set 01M6J965 serves intent m4-runs-hold-up. It covers non-plan runs only, ship-change in practice, and leaves plan runs to draft change set 01M6H5VQ (rhiza run resume already continues a plan run whose record landed, per plan-run-resumes-landed-record). A non-plan run whose change set has landed can be resumed, by rhiza run resume and by the daemon's Resumable without a pane, when its journal shows it reached the land step (land started or passed). Its remaining steps are then land and the steps after it in ship-change.toml (wait-landed through rebase-merge, merge, wait-main, main-green, close-out, ship), and each one either skips work already done or only checks state. A landed change set whose run never reached land, such as a plain build or review run, is still refused with the existing 'no longer accepted' conflict and next: rhiza run abandon. Rejected, withdrawn and superseded change sets are refused as before

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-0a9726dd-1](../question/question-0a9726dd-1.md)
- derived-from [landed-run-resume-plan](../source/landed-run-resume-plan.md)

### In

- decided-by [run-resume-landed-past-land](../capability/rhiza/workflows/workflow-runs/run-resume-landed-past-land.md)
- decided-by [daemon-resumes-landed-past-land](../capability/rhiza/orchestration/daemon/daemon-resumes-landed-past-land.md)
- decided-by [inbox-omits-resumable-landed-run](../capability/rhiza/orchestration/inbox/inbox-omits-resumable-landed-run.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
