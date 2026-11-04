# narrow-unaccepted-run-in-inbox

decision · `01M6JARXXMYR43XN6S5ATVP6GA`

## reasoning

The change set's intent is to let post-land runs on landed change sets be resumed, and the plan already narrowed run-resume, manual-resume-refuses-unaccepted and daemon-resumes-accepted-only for that, with daemon-resumes-accepted-only pointing at the exception by handle. unaccepted-run-in-inbox is the one sibling left saying every interrupted run on a non-accepted change set waits in the inbox, which contradicts inbox-omits-resumable-landed-run case (1). Its existing test uses journals with no step entries, so it still holds once the claim is limited to runs that never reached land, and inbox-omits-resumable-landed-run case (3) also covers that. Narrowing the statement and acceptance, with the exception named by handle as in daemon-resumes-accepted-only, settles the contradiction without changing behaviour or dropping the test. Keeping the node unchanged would leave the model contradicting itself, and retiring it would lose the inbox claim for runs that never landed

## statement

This change set narrows unaccepted-run-in-inbox the way it narrows its siblings. Its statement becomes: an interrupted run on a change set no longer accepted waits in the inbox, naming the change set's state and rhiza run abandon as the next step, except a non-plan run on a landed change set whose journal reached its land step, which the daemon resumes (inbox-omits-resumable-landed-run, daemon-resumes-landed-past-land). Its acceptance becomes: a test with an interrupted run on a landed change set whose journal never reached land lists it in rhiza inbox with the state and the abandon command, and the daemon does not resume it. The existing daemon-resume.txtar cases (LANDED00001/LANDEDAUTO1, journals with no step entries) stay as the test for that narrowed claim

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6j962-r3-f1](../question/review-01m6j962-r3-f1.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
