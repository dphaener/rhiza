# daemon-resume-once

capability · `01M61R7MP3HGEB60GJV4M3CWQA` · parent [daemon](../daemon.md)

## acceptance

A testscript with a fixture journal that is unfinished and names a dead process runs rhiza daemon --once twice with a fake Rex: the first starts one detached rhiza run resume --foreground on that run and journals it, the second (with the run interrupted again) starts none; a failed run's journal gets no resume.

## statement

The daemon resumes an interrupted run whose change set is accepted, or whose change set has landed when it is a plan-workflow run (daemon-resumes-landed-plan-run) or a non-plan run whose journal reached its land step (daemon-resumes-landed-past-land), once, as a detached rhiza run resume --foreground <run> process outside Rex (daemon-headless-run-drivers), and journals that resume in the run's journal; it leaves alone every other interrupted run whose change set is not accepted. A run interrupted again after that resume, and every failed run, gets no second plain resume: it is left to the watchdog's policy (daemon-watchdog and its leaves), which may resume it with --retry within the per-run retry budget; only when that policy does not apply, is turned off in rhiza.toml, or its budget is spent does the run stay listed in the inbox as a run the daemon gave up on, or as failed.

## Links

### Out

- depends-on [run-process-record](../../workflows/workflow-runs/run-process-record.md)
- depends-on daemon-run-tabs (retired)
- decided-by [answer-8cd68bc7-15](../../../../decision/answer-8cd68bc7-15.md)
- decided-by [decision-plan-landed-resumable](../../../../decision/decision-plan-landed-resumable.md)
- decided-by [daemon-resume-once-defers-to-watchdog](../../../../decision/daemon-resume-once-defers-to-watchdog.md)

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
- verifies [pr-375-daemon-watchdog](../../../../evidence/pr-375-daemon-watchdog.md)

## History

- `01M61R3R8BS63ZGH0BVS7QQ5SS` The daemon: rhiza daemon starts ship-change runs on accepted, ready change sets that serve a milestone and have leaves, in roadmap order under a concurrency cap, each in its own Rex tab, resumes an interrupted run once, and reports through rhiza daemon status.
- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
- `01M6FENEJE4B11EBYXDGP9T052` Run drivers run headless as detached children of the daemon, every run starts through the daemon, and agent steps open their own Rex tabs
- `01M6J965BRJ2ZTDT7Y311PM96G` Fix A run interrupted after its change set landed can never be resumed: rhiza run resume and the daemon's resume both require the change set to be accepted (run.go resume guard, daemon Resumable's accepted check), so a ship-change run stopped at wait-main or close-out (or a plan run after its record landed) is stranded and must be abandoned and closed out by hand. Allow resuming a run whose change set landed when its remaining steps are post-land ones (wait-main, close-out and later), and let the daemon resume those headless
- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
- `01M6WHRFPM0C0X3H6VBFT4R19Y` Daemon watchdog: the daemon notices work that has stopped progressing and recovers it without the owner or an outside session watching. Each poll it checks for: a run with no journal or session progress for longer than its step's limit (per step kind, defaults a workflow can override), a run failed on a transient signature (vanished go $WORK, git index.lock, network/GitHub 5xx, CI flake that passes on rerun), a request stuck planning or failed, a green PR its run never merged, a PR with no CI, and its own poll falling behind. It recovers by the supported path (resume --retry with a per-run retry budget, rerun failed CI once, rebase), logs one line per action, and after the budget is spent leaves an inbox item naming what it tried. For every recovery caused by a bug it files one bug request with the evidence (run and request IDs, step, log lines), deduplicated against open requests. rhiza health shows the watchdog's last actions. All of it is policy with shipped defaults (decisions workflows-are-the-rules and warn-dont-safeguard): limits, retry budgets and whether to file bugs live in rhiza.toml. Today an outside Claude session loops over health, inbox, request ls and gh pr list to do this by hand (watch log in the orchestrator's scratchpad, 2026-11-08).
