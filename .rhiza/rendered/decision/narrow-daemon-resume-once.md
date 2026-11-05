# narrow-daemon-resume-once

decision · `01M6JAVFHDFYDB9KZGM4FDEQJY`

## reasoning

The change set's intent (m4-runs-hold-up, landed-run-resume-scope) is that the daemon resumes a non-plan run on a landed change set once it reached land. The plan already narrowed daemon-resumes-accepted-only, manual-resume-refuses-unaccepted and run-resume for that, naming the exception by handle, and narrow-unaccepted-run-in-inbox settled the sibling finding f1 the same way. daemon-resume-once is the last node still saying every interrupted run on a non-accepted change set is left alone, which contradicts daemon-resumes-landed-past-land. Its once-only rule, gave-up and failed handling all still hold for the landed case (daemon-resumes-landed-past-land says 'once and headless', and inbox-omits-resumable-landed-run lists such a run as gave-up or failed), so only the not-accepted clause needs the exception. Its acceptance tests an accepted run's single resume and a failed run's none, which the change does not touch, so it stays. Narrowing keeps the model consistent without changing behaviour beyond what the plan already decided or dropping a test

## statement

This change set narrows daemon-resume-once the way it narrowed daemon-resumes-accepted-only. Its statement becomes: the daemon resumes an interrupted run whose change set is accepted, or a non-plan run on a landed change set whose journal reached its land step (daemon-resumes-landed-past-land), once, as a detached rhiza run resume --foreground <run> process outside Rex (daemon-headless-run-drivers), and journals that resume in the run's journal; it leaves alone an interrupted run whose change set is not accepted (outside that exception), a run interrupted again after that resume, which the inbox lists as a run the daemon gave up on, and every failed run, which the inbox lists as failed. Its acceptance stays as it is

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6j962-r3-f2](../question/review-01m6j962-r3-f2.md)

## History

- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
