# daemon-resumes-accepted-only-acceptance-split

decision · `01M6JAXW1MCBKSCB2K325PYJ2J`

## reasoning

The node's statement now carries the landed-past-land exception (daemon-resumes-landed-past-land) but its acceptance still says landed runs stay unresumed, so the two fields of one node disagree. The change set's intent is exactly to let the daemon resume landed runs whose remaining steps are post-land, so the statement is right and the acceptance is stale. The existing fixtures never reach land, so qualifying the landed case keeps them passing while removing the reading that every landed run must stay unresumed; the resumed landed case is the new behaviour the plan adds and belongs in the acceptance too.

## statement

Qualify daemon-resumes-accepted-only's acceptance to match its statement: a test resumes an interrupted run on an accepted change set, and one on a landed change set whose non-plan run reached its land step; it leaves interrupted runs unresumed on landed change sets that never reached land, and on rejected, withdrawn, proposed and draft change sets

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6j962-r3-f3](../question/review-01m6j962-r3-f3.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
