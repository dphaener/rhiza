# daemon-resume-once

capability · `01M61R7MP3HGEB60GJV4M3CWQA` · parent [daemon](../daemon.md)

## acceptance

A testscript with a fixture journal that is unfinished and names a dead process runs rhiza daemon --once twice with a fake Rex: the first starts one detached rhiza run resume --foreground on that run and journals it, the second (with the run interrupted again) starts none; a failed run's journal gets no resume.

## statement

The daemon resumes an interrupted run whose change set is accepted, or whose change set has landed when it is a plan-workflow run (daemon-resumes-landed-plan-run) or a non-plan run whose journal reached its land step (daemon-resumes-landed-past-land), once, as a detached rhiza run resume --foreground <run> process outside Rex (daemon-headless-run-drivers), and journals that resume in the run's journal; it leaves alone every other interrupted run whose change set is not accepted, a run interrupted again after that resume, which the inbox lists as a run the daemon gave up on, and every failed run, which the inbox lists as failed.

## Links

### Out

- depends-on [run-process-record](../../workflows/workflow-runs/run-process-record.md)
- depends-on daemon-run-tabs (retired)
- decided-by [answer-8cd68bc7-15](../../../../decision/answer-8cd68bc7-15.md)
- decided-by [decision-plan-landed-resumable](../../../../decision/decision-plan-landed-resumable.md)

### In

- depends-on [inbox-list](../inbox/inbox-list.md)
- depends-on [daemon-resumes-accepted-only](daemon-resumes-accepted-only.md)
- depends-on [daemon-resumes-landed-past-land](daemon-resumes-landed-past-land.md)
- blocks [review-01m6j962-r3-f2](../../../../question/review-01m6j962-r3-f2.md)
- verifies [daemon-resume-once-verified](../../../../evidence/daemon-resume-once-verified.md)
- verifies [daemon-resume-once-qa-1](../../../../evidence/daemon-resume-once-qa-1.md)
- verifies [daemon-resume-once-verified-r1](../../../../evidence/daemon-resume-once-verified-r1.md)
- verifies [daemon-resume-once-qa-3](../../../../evidence/daemon-resume-once-qa-3.md)
- verifies [pr-52-the-daemon](../../../../evidence/pr-52-the-daemon.md)
- verifies [pr-64-resume-only-in-flight-runs](../../../../evidence/pr-64-resume-only-in-flight-runs.md)
- verifies [pr-210-run-drivers-run-headless-as](../../../../evidence/pr-210-run-drivers-run-headless-as.md)
- verifies [pr-246-fix-a-run-interrupted-after](../../../../evidence/pr-246-fix-a-run-interrupted-after.md)
- verifies [pr-276-fix-inbox-and-daemon-treat-an](../../../../evidence/pr-276-fix-inbox-and-daemon-treat-an.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
