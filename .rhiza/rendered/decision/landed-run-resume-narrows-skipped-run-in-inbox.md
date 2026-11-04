# landed-run-resume-narrows-skipped-run-in-inbox

decision · `01M6JB0G1QSTZ3R9TGP9ABGECS`

## reasoning

skipped-run-in-inbox came from the plan interview resume-only-in-flight-plan-interview, where the owner accepted the agent's recommended 'Show in inbox' option; the question then was how to surface runs the daemon cannot safely continue, not whether a run past land should stay stranded. The owner's later intent m4-runs-hold-up says a ship-change run finishes without hand repair, and this change set's intent is exactly to stop abandoning and closing out by hand a run stopped at wait-main or close-out. A run that reached land has only idempotent or check-only steps left (ship land reports landed=false on an already-landed change set), so resuming it carries none of the risk the inbox ruling guards against. The model already narrows the sibling nodes the same way, naming the exception by handle (daemon-resumes-accepted-only, narrow-unaccepted-run-in-inbox, narrow-daemon-resume-once), so naming skipped-run-in-inbox explicitly makes the change visible and consistent without escalating. Escalating would block the run on a ruling the owner's own intent already answers, and leaving it silent is what finding f4 rightly flagged

## statement

landed-run-resume-scope knowingly narrows the owner decision skipped-run-in-inbox; it is settled here, not sent back to the owner. skipped-run-in-inbox keeps holding for every interrupted run whose change set is no longer accepted (rejected, withdrawn, superseded, and landed runs whose journal never reached land): not resumed, waits in the inbox with the change set's state and rhiza run abandon as next. The one exception is a non-plan run on a landed change set whose journal shows it reached its land step: it is resumed by rhiza run resume and by the daemon (landed-run-resume-scope, daemon-resumes-landed-past-land) and is not listed as skipped (inbox-omits-resumable-landed-run). Readers of skipped-run-in-inbox should read it with this exception

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6j962-r3-f4](../question/review-01m6j962-r3-f4.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
